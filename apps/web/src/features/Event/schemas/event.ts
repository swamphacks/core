import { z } from "zod";

export const EventSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  description: z.string().nullable(),
  location: z.string().nullable(),
  location_url: z.url().nullable(),
  max_attendees: z.number().int().nullable(),
  application_open: z.coerce.date(),
  application_close: z.coerce.date(),
  rsvp_deadline: z.coerce.date().nullable(),
  decision_release: z.coerce.date().nullable(),
  start_time: z.coerce.date(),
  end_time: z.coerce.date(),
  website_url: z.url().nullable(),
  is_published: z.boolean(),
  created_at: z.coerce.date().nullable(),
  updated_at: z.coerce.date().nullable(),
});

export type Event = z.infer<typeof EventSchema>;
