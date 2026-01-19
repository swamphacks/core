export const Intent = {
  IDENT: "INDET",
} as const;

export type Intent = (typeof Intent)[keyof typeof Intent];
