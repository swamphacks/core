import type { OAuthProvider } from "../types/types";

// https://www.typescriptlang.org/docs/handbook/release-notes/typescript-5-0.html#const-type-parameters
// In short, this allows us to infer the narrowest type possible for each provider object without using "as const" on them
export const createProvider = <const T extends OAuthProvider>(provider: T) => {
  return provider;
};
