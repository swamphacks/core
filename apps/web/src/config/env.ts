import { z } from "zod";

const isRuntime =
  typeof window !== "undefined" &&
  typeof (window as any).ENV !== "undefined" &&
  import.meta.env.PROD;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const runtimeEnv = (window as any).ENV || {};

const rawEnvSource = isRuntime ? runtimeEnv : import.meta.env;

const envSchema = z.object({
  BASE_API_URL: z.url(),
  DISCORD_OAUTH_CLIENT_ID: z.string(),
  VITE_ALLOWED_HOSTS: z.array(z.string()).optional(),
});

export type Env = z.infer<typeof envSchema>;

const getEnvValue = (
  runtimeKey: string,
  devKey: string,
): string | undefined => {
  if (isRuntime) {
    return rawEnvSource[runtimeKey];
  }
  return rawEnvSource[devKey];
};

const getAllowedHosts = () => {
  const hostsString = getEnvValue("VITE_ALLOWED_HOSTS", "VITE_ALLOWED_HOSTS");

  if (hostsString) {
    try {
      return JSON.parse(hostsString);
    } catch (e) {
      console.error("Failed to parse VITE_ALLOWED_HOSTS as JSON array:", e);
      return undefined;
    }
  }
  return undefined;
};

export const env = envSchema.parse({
  BASE_API_URL: getEnvValue("VITE_BASE_API_URL", "VITE_BASE_API_URL"),
  DISCORD_OAUTH_CLIENT_ID: getEnvValue(
    "VITE_DISCORD_OAUTH_CLIENT_ID",
    "VITE_DISCORD_OAUTH_CLIENT_ID",
  ),
  VITE_ALLOWED_HOSTS: getAllowedHosts(),
});
