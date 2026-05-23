import Auth from "./auth";
import { authConfig } from "./auth/config";
import { Discord } from "./auth/providers";
import { queryClient } from "./tanstack-query-client";
import { useUserQueryKey } from "./auth/hooks/useUser";

export const auth = Auth({
  providers: [Discord],
  redirectUri: authConfig.OAUTH_REDIRECT_URL,

  hooks: {
    afterLogout: async () => {
      await queryClient.invalidateQueries({ queryKey: useUserQueryKey });
    },
  },
});
