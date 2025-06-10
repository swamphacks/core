import Auth from "./auth";
import { authConfig } from "./auth/config";
import { Discord } from "./auth/providers";

export const authClient = Auth({
  providers: [Discord],
  redirectUri: authConfig.OAUTH_REDIRECT_URL,
});
