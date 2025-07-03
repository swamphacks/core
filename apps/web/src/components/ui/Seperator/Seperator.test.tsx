import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Separator } from ".";

describe("Separator component", () => {
  it("renders with correct title", () => {
    render(<Separator />);
    const separator = screen.getByTestId("separator");
    expect(separator).toBeInTheDocument();
  });
});
