import { useUser } from "./hooks/useUser";
import { _oauthSignIn } from "./services/oauth";
import { _getUser, _internalLogOut } from "./services/session";
import type { AuthConfig } from "./types";

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    oauth: {
      signIn: _oauthSignIn(config),
    },

    // Generic functions
    logOut: _internalLogOut,
    getUser: _getUser,
    useUser,
    getSession: () => {}, // This needs backend implementation to finish
  };
}
