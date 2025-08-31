import { z } from "zod";

const envSchema = z.object({
  BASE_API_URL: z.string().url(),
  DISCORD_OAUTH_CLIENT_ID: z.string(),
  VITE_ALLOWED_HOSTS: z.array(z.string()).optional(),
});

export type Env = z.infer<typeof envSchema>;

export const env = envSchema.parse({
  BASE_API_URL: import.meta.env.VITE_BASE_API_URL,
  DISCORD_OAUTH_CLIENT_ID: import.meta.env.VITE_DISCORD_OAUTH_CLIENT_ID,
  VITE_ALLOWED_HOSTS: import.meta.env.VITE_ALLOWED_HOSTS
    ? JSON.parse(import.meta.env.VITE_ALLOWED_HOSTS)
    : undefined,
});
