import * as z from "zod";
import { Intent } from "./intent";

// Used for checking attendees into an event
export type IdentifyIntent = {
  intent: typeof Intent.IDENT;
  user_id: string;
};

export type QRIntent = IdentifyIntent;
export type IntentParseError = "MALFORMED_INTENT_HEADER" | "MALFORMED_BODY";

// Can be abstracted away later
export type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

export function parseQrIntent(
  input: string,
): Result<QRIntent, IntentParseError> {
  if (input.trim().length <= 0) {
    return {
      ok: false,
      error: "MALFORMED_INTENT_HEADER",
    };
  }

  const [head, body] = input.trim().split("::");

  switch (head.toUpperCase()) {
    case "IDENT": {
      const res = parseCheckIn(body);

      if (!res.ok) {
        return {
          ok: false,
          error: "MALFORMED_BODY",
        };
      }

      return {
        ok: true,
        value: {
          intent: Intent.IDENT,
          ...res.value,
        },
      };
    }
    default: {
      console.log(`Header not recognized : ${head.toUpperCase()}`);
      return {
        ok: false,
        error: "MALFORMED_INTENT_HEADER",
      };
    }
  }
}

const CheckInSchema = z.object({
  user_id: z.uuid(),
});

type CheckInFields = z.infer<typeof CheckInSchema>;

type CheckInParseError = "CHECK_IN_PARSE_ERROR";

function parseCheckIn(input: string): Result<CheckInFields, CheckInParseError> {
  const { success, data } = CheckInSchema.safeParse({
    user_id: input,
  });

  if (!success) {
    return {
      ok: false,
      error: "CHECK_IN_PARSE_ERROR",
    };
  }

  return {
    ok: true,
    value: data,
  };
}
