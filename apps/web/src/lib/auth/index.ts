import Cookies from "js-cookie";
import type { AuthConfig, OAuthState, User } from "./types";
import { createOAuthRequestParams, isInSession } from "./utils";
import { authConfig } from "./config";
import { HTTPError } from "ky";
import { api } from "../api";
import { useQuery } from "@tanstack/react-query";

const userApi = api.extend((options) => ({
  prefixUrl: `${options.prefixUrl}/protected`,
}));

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

  return (providerId: ProviderId, redirectTo?: string) => {
    const provider = config.providers.find(
      (provider) => provider.id === providerId,
    );

    if (!provider) {
      throw new Error(`OAuth provider "${providerId}" not found.`);
    }

    // Validate redirectTo url
    if (redirectTo) {
      try {
        new URL(redirectTo);
      } catch {
        throw new Error(`returnTo URL "${redirectTo}" is invalid.`);
      }
    }

    const url = new URL(provider.authorization.url);

    // Set state parameter to prevent CSRF attacks
    const state: OAuthState = {
      nonce: crypto.randomUUID(),
      provider: providerId,
      redirect: redirectTo,
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

/**
 * Gets the current user details if they are authenticated. If a session cookie is found, this function performs a
 * raw network request, so the returned data is up to date. Results are not cached.
 *
 * @returns An object containing the user data and error (if any). If there is an error, user is `null`.
 */
async function getUser() {
  try {
    if (!isInSession()) {
      throw new Error("cookie session not found");
    }

    const user = await userApi(authConfig.AUTH_ME_URL, {
      credentials: "include",
    }).json<User>();

    return { user, error: null };
  } catch (error) {
    // TODO: extract to error handling function
    let errorMsg = "auth error";
    if (error instanceof HTTPError) {
      try {
        const errorResponse = await error.response.json();
        errorMsg = errorResponse.message || errorMsg;
        console.error(errorResponse.message);
      } catch {
        console.error("Failed to parse error response");
      }
    } else {
      // console.error(`Network or system error: ${error.message}`);
    }

    return { user: null, error: errorMsg };
  }
}

/**
 * A React hook that retrieves current user details if they are authenticated. It uses
 * Tanstack Query to fetch data using the `getUser` function as the query function.
 * Results will be cached for 5 minutes.
 *
 * @returns An object containing the user data, error (if any) and a promise. If there is an error, user is `null`. The promise can be awaited,
 * and will resolve with the user data.
 */
function useUser() {
  const result = useQuery({
    queryKey: ["user"],
    queryFn: async () => {
      const { user, error } = await getUser();
      if (error) {
        throw new Error(error);
      }

      // user is guaranteed to not be null here
      return user as User;
    },
    staleTime: 5 * 60 * 1000, // cache for 5 minutes,
    retry: (failureCount) => {
      // console.log(error);
      return isInSession() && failureCount < 3;
    },
  });

  return {
    user: result.data,
    error: result.error,
    // https://tanstack.com/query/latest/docs/framework/react/reference/useQuery
    // this is an experimental feature so I'm not sure if it will cause problem later or not
    promise: result.promise,
  };
}

function logOut() {
  throw new Error("Not implemented");
}

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    oauth: {
      signIn: _oauthSignIn(config),
    },

    // Generic functions
    logOut,
    getUser,

    // Hooks
    useUser,
  };
}
