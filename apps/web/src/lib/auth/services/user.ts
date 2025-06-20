import { api } from "@/lib/ky";
import * as ky from "ky";
import type { AuthUserResponse } from "../types/auth";
import { errorSchema } from "../types/error";
import { userContextSchema } from "../types/user";
import { authConfig } from "../config";
import type { AuthConfig } from "../types";

export function _logout(config: AuthConfig) {
  return async () => {
    await config.hooks.beforeLogout?.();

    const res = await api.post(authConfig.AUTH_LOGOUT_URL);

    await config.hooks.afterLogout?.();

    if (!res.ok) {
      throw new Error("Failed to log out");
    }

    // Client will handle redirect
    window.location.href = "/";
  };
}

export async function _getUser(): Promise<AuthUserResponse> {
  try {
    console.log("fetching user info...");
    const res = await api
      .get(authConfig.AUTH_ME_URL, {
        retry: 0, // disable retry
      })
      .json();

    // Attempt to parse response in AuthUserResponse schema
    const userContext = userContextSchema.safeParse(res);
    if (!userContext.success) throw userContext.error;

    console.log("fetched user: ", userContext.data);

    return {
      user: userContext.data,
      error: null,
    };
  } catch (error) {
    return handleError(error);
  }
}

async function handleError(error: unknown) {
  if (error instanceof ky.HTTPError && error.response.status === 401) {
    // Not logged in, return null and null
    return {
      user: null,
      error: null,
    };
  } else if (error instanceof ky.HTTPError) {
    // Attempt to parse body in ErrorResponse schema
    const errorResponse = errorSchema.safeParse(await error.response.json());
    if (!errorResponse.success) {
      // Almost certainly an internal error as that is what we return from API when our response fails
      // Questions? Ask @alexwala on discord
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
  } else {
    // Handle other types of errors (parse, network, etc.)
    return {
      user: null,
      error: {
        error: "unknown_err",
        message: "An unknown error occured, please try again later",
      },
    };
  }
}
