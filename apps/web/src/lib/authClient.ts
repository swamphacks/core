import Auth from "./auth";
import { authConfig } from "./auth/config";
import { Discord } from "./auth/providers";
import { queryClient } from "./query";
import { queryKey as useUserQueryKey } from "./auth/hooks/useUser";

export const auth = Auth({
  providers: [Discord],
  redirectUri: authConfig.OAUTH_REDIRECT_URL,

  hooks: {
    afterLogout: async () => {
      await queryClient.invalidateQueries({ queryKey: useUserQueryKey });
    },
  },
});
