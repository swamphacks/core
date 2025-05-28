import Cookies from "js-cookie";
import { useEffect, useRef } from "react";
import { APP_URL } from "@/utils/url";
import type { User } from "@/features/common/Auth/types";

export const useAuth = () => {
  // // https://github.com/TanStack/router/discussions/1668#discussioncomment-10634735
  const authStateRef = useRef<PromiseWithResolvers<User | null>>(undefined);
  if (!authStateRef.current) {
    authStateRef.current = Promise.withResolvers<User | null>();
  }

  useEffect(() => {
    async function getAuthState() {
      const authCookie = Cookies.get("auth");

      // If auth cookie isn't set, user is unauthenticated, so no need to send a verify request
      if (!authCookie) {
        authStateRef.current!.resolve(null);
        return;
      }

      try {
        const res = await fetch(APP_URL.VERIFY_AUTH, {
          credentials: "include",
        });

        authStateRef.current!.resolve(await res.json());
      } catch {
        // TODO: Add proper error handling and alert the user
        authStateRef.current!.resolve(null);
      }
    }

    getAuthState();
  }, []);

  return authStateRef.current;
};
