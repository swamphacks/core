import Cookies from "js-cookie";
import type {
  OauthProvider,
  BaseOAuthParams,
  AuthConfig,
  OauthState,
} from "./types";

const buildOauthParams = (
  provider: OauthProvider,
  redirectUri: string,
  state: OauthState,
): BaseOAuthParams => ({
  response_type: "code",
  scope: provider.authorization.scopes,
  client_id: provider.authorization.clientId,
  redirect_uri: redirectUri,
  state: btoa(JSON.stringify(state)),
});

/*
 * This internal function enables strict typing for "auth.signIn()".
 * It ensures that the "providerId" argument can only be one of the "id" values
 * from the providers specified in the Auth configuration, rather than any string.
 */
function _signInInternal<T extends AuthConfig>(config: T) {
  type ProviderId = T["providers"][number]["id"];

  return (providerId: ProviderId, returnTo?: string) => {
    const provider = config.providers.find(
      (provider) => provider.id === providerId,
    );

    if (!provider) {
      throw new Error(`OAuth provider "${providerId}" not found.`);
    }

    // Validate returnTo url
    if (returnTo) {
      try {
        new URL(returnTo);
      } catch {
        throw new Error(`returnTo URL "${returnTo}" is invalid.`);
      }
    }

    const url = new URL(provider.authorization.url);

    // Set state parameter to prevent CSRF attacks
    const state: OauthState = {
      nonce: crypto.randomUUID(),
      provider: providerId,
      redirect: returnTo,
    };

    Cookies.set("sh_auth_nonce", state.nonce, {
      path: "/",
      // Set sameSite to lax so the cookie can be sent when the authorization server redirects to redirect_uri
      sameSite: "lax",
      secure: true,
    });

    const params = buildOauthParams(provider, config.redirect_uri, state);

    for (const [param, value] of Object.entries(params)) {
      url.searchParams.set(param, value);
    }

    window.location.href = url.toString();
  };
}

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    signIn: _signInInternal(config),
    signOut: () => {
      // TODO: implement logout
    },
  };
}
