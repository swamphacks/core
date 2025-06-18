import { api } from "@/lib/ky";
import * as ky from "ky";
import type { AuthMeResponse } from "../types/auth";
import { errorSchema } from "../types/error";
import { userContextSchema } from "../types/user";

export async function _internalLogOut(): Promise<void> {
  const res = await api.post("auth/logout");
  if (!res.ok) {
    throw new Error("Failed to log out");
  }

  // Client will handle redirect
}

export async function _getUser(): Promise<AuthMeResponse> {
  try {
    const res = await api.get("auth/me").json();

    // Attempt to parse response in AuthMeResponse schema
    const userContext = userContextSchema.safeParse(res);
    if (!userContext.success) throw userContext.error;

    return {
      user: userContext.data,
      error: null,
    };
  } catch (error) {
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
}
