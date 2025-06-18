import Auth from "./auth";
import { authConfig } from "./auth/config";
import { Discord } from "./auth/providers";
import { queryClient } from "./query";

export const auth = Auth({
  providers: [Discord],
  redirectUri: authConfig.OAUTH_REDIRECT_URL,
});

export async function logout() {
  try {
    await auth.logOut();
    await queryClient.invalidateQueries({ queryKey: ["auth", "user"] });
  } catch (error) {
    console.error("Logout failed:", error);
  }
}
