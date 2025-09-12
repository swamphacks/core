import { z } from "zod";

export const userContextSchema = z.object({
  userId: z.string().uuid(),
  email: z.email(),
  preferredEmail: z.preprocess(
    (val: string) => (val === "" ? undefined : val),
    z.email().optional(),
  ),
  name: z.string(),
  onboarded: z.boolean(),
  image: z.string().nullable().optional(),
  role: z.enum(["user", "superuser"]),
  emailConsent: z.boolean(),
});

export type UserContext = z.infer<typeof userContextSchema>;
