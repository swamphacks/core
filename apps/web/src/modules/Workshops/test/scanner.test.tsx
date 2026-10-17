import {
  act,
  fireEvent,
  render,
  screen,
  cleanup,
} from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import QrWorkshopScanner from "../QrWorkshopScanner";

const mocks = vi.hoisted(() => ({
  post: vi.fn(),
  scan: undefined as
    | undefined
    | ((codes: { rawValue: string }[]) => Promise<void>),
}));
vi.mock("@/lib/ky", () => ({ api: { post: mocks.post } }));
vi.mock("@tanstack/react-query", () => ({
  useQuery: () => ({
    data: [
      {
        id: "session-1",
        title: "Build a gator",
        type: "workshop",
        startTime: "2026-10-17T14:00:00Z",
      },
    ],
    isError: false,
  }),
}));
vi.mock("@/components/ui/Modal", () => ({
  Modal: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));
vi.mock("@yudiel/react-qr-scanner", () => ({
  Scanner: ({ onScan }: { onScan: typeof mocks.scan }) => {
    mocks.scan = onScan;
    return <div>Camera</div>;
  },
}));
vi.mock("~icons/tabler/circle-check-filled", () => ({ default: () => null }));

beforeEach(() => {
  mocks.post.mockReset();
  mocks.scan = undefined;
});
afterEach(cleanup);

function selectSession() {
  render(<QrWorkshopScanner isOpen onClose={() => {}} />);
  fireEvent.change(screen.getByRole("combobox"), {
    target: { value: "session-1" },
  });
}
const userId = "6603e995-7ac7-4af9-903b-ea5d4b0559b2";

describe("Workshop QR scanning", () => {
  it("records the selected workshop and suppresses repeated camera reads", async () => {
    mocks.post.mockReturnValue({
      json: async () => ({
        name: "Luca",
        workshopTitle: "Build a gator",
        alreadyRecorded: false,
      }),
    });
    selectSession();
    await act(async () => {
      await mocks.scan?.([{ rawValue: `IDENT::${userId}` }]);
    });
    expect(mocks.post).toHaveBeenCalledWith("workshops/session-1/attendance", {
      json: { userId },
      retry: 0,
    });
    expect(screen.getByText("Attendance recorded.")).toBeDefined();
    await act(async () => {
      await mocks.scan?.([{ rawValue: `IDENT::${userId}` }]);
    });
    expect(mocks.post).toHaveBeenCalledTimes(1);
  });
  it("shows an existing NFC/QR record without claiming another attendance", async () => {
    mocks.post.mockReturnValue({
      json: async () => ({
        name: "Luca",
        workshopTitle: "Build a gator",
        alreadyRecorded: true,
      }),
    });
    selectSession();
    await act(async () => {
      await mocks.scan?.([{ rawValue: `IDENT::${userId}` }]);
    });
    expect(
      screen.getByText("Already recorded — no duplicate added."),
    ).toBeDefined();
  });
  it("rejects non-hacker codes without an API write", async () => {
    selectSession();
    await act(async () => {
      await mocks.scan?.([{ rawValue: "https://example.com" }]);
    });
    expect(mocks.post).not.toHaveBeenCalled();
    expect(
      screen.getByText("This is not a SwampHacks hacker QR code."),
    ).toBeDefined();
  });
});
