import { createProvider } from "./createProvider";
import { authConfig } from "../config";

export const Discord = createProvider({
  id: "discord",
  authorization: {
    url: authConfig.DISCORD_OAUTH_BASE_URL,
    scopes: "identify email",
    clientId: authConfig.DISCORD_OAUTH_CLIENT_ID,
  },
});
