import type { BaseOAuthParams, OAuthProvider, OAuthState } from "./types";
import Cookies from "js-cookie";

export const createOAuthRequestParams = (
  provider: OAuthProvider,
  redirectUri: string,
  state: OAuthState,
): BaseOAuthParams => ({
  response_type: "code",
  scope: provider.authorization.scopes,
  client_id: provider.authorization.clientId,
  redirect_uri: redirectUri,
  state: btoa(JSON.stringify(state)),
});

export function isInSession() {
  return !!Cookies.get("sh_session");
}
