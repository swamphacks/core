import { authConfig } from "./config";
import type { OAuthProvider } from "./types";

// https://www.typescriptlang.org/docs/handbook/release-notes/typescript-5-0.html#const-type-parameters
// In short, this allows us to infer the narrowest type possible for each provider object without using "as const" on them
export const createProvider = <const T extends OAuthProvider>(provider: T) => {
  return provider;
};

export const Discord = createProvider({
  id: "discord",
  authorization: {
    url: authConfig.DISCORD_OAUTH_BASE_URL,
    scopes: "identify email",
    clientId: authConfig.DISCORD_OAUTH_CLIENT_ID,
  },
});

// Add more providers here as needed
