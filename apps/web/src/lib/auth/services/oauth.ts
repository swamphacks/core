import Cookies from "js-cookie";
import type { AuthConfig, OAuthState } from "../types/types";
import { createOAuthRequestParams } from "../utils/oauth";

/**
 * Internal function to handle OAuth sign-in for configured providers.
 *
 * @template T - The AuthConfig type containing OAuth providers and configuration.
 * @param {T} config - The authentication configuration object.
 * @returns A function that initiates OAuth sign-in for the specified provider ID and optional return URL.
 *
 * @throws {Error} Throws if the provider ID is not found in the config.
 */
export function _oauthSignIn<T extends AuthConfig>(config: T) {
  type ProviderId = T["providers"][number]["id"];

  return async (providerId: ProviderId, redirect?: string) => {
    const provider = config.providers.find(
      (provider) => provider.id === providerId,
    );

    if (!provider) {
      throw new Error(`OAuth provider "${providerId}" not found.`);
    }

    await config.hooks.beforeLogin?.();

    const url = new URL(provider.authorization.url);

    // Set state parameter to prevent CSRF attacks
    const state: OAuthState = {
      nonce: crypto.randomUUID(),
      provider: providerId,
      redirect,
    };

    Cookies.set("sh_auth_nonce", state.nonce, {
      path: "/",
      sameSite: "lax",
      secure: import.meta.env.DEV ? false : true,
      //TODO: change this to an env variable
      domain: import.meta.env.DEV ? "localhost" : ".swamphacks.com",
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
