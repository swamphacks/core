import { z } from "zod";

export const userContextSchema = z.object({
  userId: z.string().uuid(),
  name: z.string(),
  onboarded: z.boolean(),
  image: z.string().nullable().optional(),
  role: z.enum(["user", "superuser"]),
});

export type UserContext = z.infer<typeof userContextSchema>;
