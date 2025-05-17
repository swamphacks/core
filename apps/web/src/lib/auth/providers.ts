import type { OauthProvider } from "./types";

// https://www.typescriptlang.org/docs/handbook/release-notes/typescript-5-0.html#const-type-parameters
// In short, this allows us to infer the narrowest type possible for each provider object without using "as const" on them
const createProvider = <const T extends OauthProvider>(provider: T) => {
  return provider;
};

export const Discord = createProvider({
  id: "discord",
  authorization: {
    url: "https://discord.com/oauth2/authorize",
    scopes: "identify email",
    clientId: import.meta.env["SWAMPHACKS_DISCORD_OAUTH_CLIENT_ID"],
  },
});
