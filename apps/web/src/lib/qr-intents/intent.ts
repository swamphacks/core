export const Intent = {
  CHECK_IN: "CHECK_IN",
  REDEEM: "REDEEM",
} as const;

export type Intent = (typeof Intent)[keyof typeof Intent];
