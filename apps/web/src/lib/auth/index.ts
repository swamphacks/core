import Cookies from "js-cookie";
import type { AuthConfig, OAuthState } from "./types";
import { createOAuthRequestParams } from "./utils";
import { z } from "zod";
import { api } from "../ky";

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

const userContextSchema = z.object({
  userId: z.string().uuid(),
  name: z.string(),
  onboarded: z.boolean(),
  image: z.string().nullable().optional(),
  role: z.enum(["user", "superuser"]),
});

type UserContext = z.infer<typeof userContextSchema>;

// No need for a separate error schema if it's just a string,
// but it's good practice if it might become more complex.
const errorSchema = z.object({
  error: z.string(),
  message: z.string(),
});

type ErrorResponse = z.infer<typeof errorSchema>;

// --- IMPROVED TYPE DEFINITION ---
// This is a cleaner, more effective discriminated union.
// 1. If 'user' is present and valid, 'error' is forbidden.
// 2. If 'user' is null, 'error' is an optional string.
type AuthMeResponse =
  | { user: UserContext; error?: never }
  | { user: null; error?: ErrorResponse };

async function _internalGetMe(): Promise<AuthMeResponse> {
  // Fetch user data from the server or authentication service.
  const res = await api.get("auth/me");

  // Check whether the response contains user data or an error.

  // If 401 or 403, return null user with error
  if (res.status === 401 || res.status === 403) {
    console.warn("Unauthorized or forbidden access to user data.");

    return {
      user: null,
    };
  }

  const json = await res.json();

  // If the response contains user data, return it
  if (userContextSchema.safeParse(json).success) {
    return {
      user: userContextSchema.parse(json),
    };
  }

  // If the response contains an error, return null user with error
  if (errorSchema.safeParse(json).success) {
    return {
      user: null,
      error: errorSchema.parse(json),
    };
  }

  return {
    user: null, // Placeholder for actual user data
    error: undefined, // Placeholder for error handling
  } satisfies AuthMeResponse;
}

async function _internalLogOut(): Promise<void> {
  const res = await api.post("auth/logout");
  if (!res.ok) {
    throw new Error("Failed to log out");
  }
}

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    oauth: {
      signIn: _oauthSignIn(config),
    },

    // Generic functions
    logOut: _internalLogOut,
    getUser: _internalGetMe,
    useSession: () => {},
  };
}
