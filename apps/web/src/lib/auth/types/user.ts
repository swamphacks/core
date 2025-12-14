import { z } from "zod";

export const userContextSchema = z.object({
  userId: z.uuid(),
  email: z.email(),
  preferredEmail: z.email().nullable().optional(),
  name: z.string(),
  onboarded: z.boolean(),
  image: z.string().nullable().optional(),
  role: z.enum(["user", "superuser"]),
  emailConsent: z.boolean(),
});

export type UserContext = z.infer<typeof userContextSchema>;
