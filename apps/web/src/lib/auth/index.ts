import Cookies from "js-cookie";
import type { AuthConfig, OAuthState } from "./types";
import { createOAuthRequestParams } from "./utils";

/**
 * Internal function to handle OAuth sign-in for configured providers.
 *
 * @template T - The AuthConfig type containing OAuth providers and configuration.
 * @param {T} config - The authentication configuration object.
 * @returns A function that initiates OAuth sign-in for the specified provider ID and optional return URL.
 *
 * @throws {Error} Throws if the provider ID is not found in the config.
 * @throws {Error} Throws if the optional `returnTo` URL is invalid.
 */
function _oauthSignIn<T extends AuthConfig>(config: T) {
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
    const state: OAuthState = {
      nonce: crypto.randomUUID(),
      provider: providerId,
      redirectUri: returnTo, // Ask about this
    };

    Cookies.set("sh_auth_nonce", state.nonce, {
      path: "/",
      // Set sameSite to lax so the cookie can be sent when the authorization server redirects to redirect_uri
      sameSite: "lax",
      secure: true,
    });

    const params = createOAuthRequestParams(
      provider,
      config.redirectUri,
      state,
    );

    for (const [param, value] of Object.entries(params)) {
      url.searchParams.set(param, value);
    }

    window.location.href = url.toString();
  };
}

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    oauth: {
      signIn: _oauthSignIn(config),
    },

    // Generic functions
    logOut: () => {},
    useUser: () => {},
    useSession: () => {},
  };
}
