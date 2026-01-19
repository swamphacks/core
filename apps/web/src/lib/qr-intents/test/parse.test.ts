import { describe, it, expect } from "vitest";
import { parseQrIntent } from "../parse.ts";
import { Intent } from "../intent.ts";

describe("Parse QR Intents", () => {
  it("parses valid check in intent", () => {
    const user_id = "aa62e8a1-b4fb-479a-8d15-e52328920d18";

    const input = `IDENT::${user_id}`;
    const result = parseQrIntent(input);

    expect(result.ok).toBe(true);

    if (!result.ok) {
      expect.fail("Returned an error when parsing QR Intents");
    }

    expect(result.value.intent).toBe(Intent.IDENT);
    expect(result.value.user_id).toBe(user_id);
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
