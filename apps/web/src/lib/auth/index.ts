import { _useUser } from "./hooks/useUser";
import { _oauthSignIn } from "./services/oauth";
import { _getUser, _logout } from "./services/user";
import type { AuthConfig } from "./types";

// The entry point to the auth library
export default function Auth<const T extends AuthConfig>(config: T) {
  return {
    oauth: {
      signIn: _oauthSignIn(config),
    },

    // Generic functions
    logOut: _logout(config),
    getUser: _getUser,
    useUser: _useUser,
  };
}
