import * as z from "zod";
import { Intent } from "./intent";

// Used for checking attendees into an event
type CheckInIntent = {
  intent: typeof Intent.CHECK_IN;
  user_id: string;
  event_id: string;
};

// Used for redeeming redeemables (food, t-shirts, etc)
type RedeemIntent = {
  intent: typeof Intent.REDEEM;
  redeemable_id: string;
  user_id: string;
  event_id: string;
};

export type QRIntent = CheckInIntent | RedeemIntent;
export type IntentParseError = "MALFORMED_INTENT_HEADER" | "MALFORMED_BODY";

// Can be abstracted away later
type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

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
    case "CHECKIN": {
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
          intent: Intent.CHECK_IN,
          ...res.value,
        },
      };
    }

    case "REDEEM": {
      console.log("REDEEM detected");
      break;
    }

    default: {
      console.log(`Header not recognized : ${head.toUpperCase()}`);
      return {
        ok: false,
        error: "MALFORMED_INTENT_HEADER",
      };
    }
  }

  return { ok: false, error: "MALFORMED_INTENT_HEADER" };
}

const CheckInSchema = z.object({
  user_id: z.uuid(),
  event_id: z.uuid(),
});

type CheckInFields = z.infer<typeof CheckInSchema>;

type CheckInParseError = "CHECK_IN_PARSE_ERROR";

function parseCheckIn(input: string): Result<CheckInFields, CheckInParseError> {
  const [user_id, event_id] = input.split("+");

  const { success, data } = CheckInSchema.safeParse({
    user_id,
    event_id,
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
