import type { AuthUserResponse } from "../types/auth";
import { errorSchema } from "../types/error";
import { userContextSchema } from "../types/user";
import { authConfig } from "../config";
import type { AuthConfig } from "../types";

export function _logout(config: AuthConfig) {
  return async () => {
    await config.hooks.beforeLogout?.();

    const res = await fetch(authConfig.AUTH_LOGOUT_URL, {
      method: "POST",
      credentials: "include",
    });

    await config.hooks.afterLogout?.();

    if (!res.ok) {
      throw new Error("Failed to log out");
    }

    // Client will handle redirect
  };
}

export async function _getUser(): Promise<AuthUserResponse> {
  console.log("fetching user info...");

  try {
    const res = await fetch(authConfig.AUTH_ME_URL, {
      credentials: "include",
    });

    if (!res.ok) {
      return await handleError(res);
    }

    // Attempt to parse response in AuthUserResponse schema
    const userContext = userContextSchema.safeParse(await res.json());
    if (!userContext.success) {
      console.error("userContext parsing failed");
      return {
        user: null,
        error: {
          error: "unknown_err",
          message: "An unknown error occured, please try again later",
        },
      };
    }

    return {
      user: userContext.data,
      error: null,
    };
  } catch (error) {
    console.log("An error occured when fetching user data", error);
    return {
      user: null,
      error: {
        error: "client_err",
        message:
          "Fetch failed likely due to network error or malformed request",
      },
    };
  }
}

async function handleError(res: Response): Promise<AuthUserResponse> {
  if (res.status === 401) {
    // Not logged in, return null and null
    return {
      user: null,
      error: null,
    };
  } else {
    // Attempt to parse body in ErrorResponse schema
    const errorResponse = errorSchema.safeParse(await res.json());
    if (!errorResponse.success) {
      // Almost certainly an internal error as that is what we return from API when our response fails
      // Questions? Ask @alexwala on discord
      console.error("error parsing failed");
      return {
        user: null,
        error: {
          error: "internal_err",
          message: "An unknown internal server error occurred",
        },
      };
    }

    return {
      user: null,
      error: errorResponse.data,
    };
  }
}
