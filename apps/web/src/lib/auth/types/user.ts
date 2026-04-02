import { z } from "zod";

export const userContextSchema = z.object({
  userId: z.uuid(),
  email: z.email(),
  preferredEmail: z.email().nullable().optional(),
  name: z.string(),
  onboarded: z.boolean(),
  image: z.string().nullable().optional(),
  role: z.enum(['admin', 'staff', 'attendee', 'applicant', 'visitor']),
  emailConsent: z.boolean(),
  checkedInAt: z.date().nullable(),
  rfid: z.string().nullable(),
});

export type UserContext = z.infer<typeof userContextSchema>;
