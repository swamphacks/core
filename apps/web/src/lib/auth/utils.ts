import type { BaseOAuthParams, OAuthProvider, OAuthState } from "./types";

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
