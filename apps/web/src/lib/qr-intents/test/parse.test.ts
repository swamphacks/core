import { describe, it, expect } from "vitest";
import { parseQrIntent } from "../parse.ts";
import { Intent } from "../intent.ts";

describe("Parse QR Intents", () => {
  it("parses valid check in intent", () => {
    const user_id = "aa62e8a1-b4fb-479a-8d15-e52328920d18";
    const event_id = "fe395fb7-7d57-49fb-90ae-d56094445e45";

    const input = `checkin::${user_id}+${event_id}`;
    const result = parseQrIntent(input);

    expect(result.ok).toBe(true);

    if (!result.ok) {
      expect.fail("Returned an error when parsing QR Intents");
    }

    expect(result.value.intent).toBe(Intent.CHECK_IN);
    expect(result.value.user_id).toBe(user_id);
    expect(result.value.event_id);
    expect(result.value.event_id).toBe(event_id);
  });

  it("parses invalid header correctly", () => {
    const input = "invalid::DoesntMatter+DoesntMatter";
    const result = parseQrIntent(input);

    expect(result.ok).toBe(false);

    if (!result.ok) {
      expect(result.error).toBe("MALFORMED_INTENT_HEADER");
    } else {
      expect.fail("Result was ok, not expected.");
    }
  });
});
