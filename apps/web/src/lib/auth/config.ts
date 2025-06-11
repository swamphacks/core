import config from "@/config";

export const authConfig = {
  // General Auth API URLs for SwampHacks Core Backend.
  OAUTH_REDIRECT_URL: "http://localhost:8080/auth/callback",
  AUTH_ME_URL: "user",
  AUTH_SESSION_URL: "session",
  AUTH_LOGOUT_URL: "logout",

  // Discord Specific OAuth Config
  DISCORD_OAUTH_BASE_URL: "https://discord.com/oauth2/authorize",
  DISCORD_OAUTH_CLIENT_ID: config.DISCORD_OAUTH_CLIENT_ID, // Copy this over from env.ts
};
