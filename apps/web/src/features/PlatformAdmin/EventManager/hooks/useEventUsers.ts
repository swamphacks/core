import { api } from "@/lib/ky";
import { EventRoleSchema, PlatformRoleSchema } from "@/lib/openapi/zodSchemas";
import { useQuery } from "@tanstack/react-query";
import { z } from "zod";

export function getEventStaffUsersQueryKey(eventId: string) {
  return ["event", eventId, "event-users"] as const;
}

export const EventUserSchema = z.object({
  id: z.uuid(), // uuid.UUID -> string with uuid validation
  name: z.string(),
  email: z.email(), // *string -> optional/nullable
  email_verified: z.boolean(),
  onboarded: z.boolean(),
  image: z.url().optional().nullable(), // *string -> optional/nullable
  created_at: z.coerce.date(),
  updated_at: z.coerce.date(),
  role: PlatformRoleSchema,
  preferred_email: z.email().optional().nullable(), // *string -> optional/nullable
  email_consent: z.boolean(),
  checked_in_at: z.coerce.date().optional().nullable(), // *time.Time -> optional/nullable
  event_role: EventRoleSchema,
});

export const GetEventUsersSchema = z.array(EventUserSchema);

export type EventUser = z.infer<typeof EventUserSchema>;

export function useEventUsers(eventId: string) {
  async function fetchEventStaffUsers(): Promise<EventUser[]> {
    const result = await api.get<EventUser[]>(`events/${eventId}/users`).json();
    return GetEventUsersSchema.parse(result);
  }

  return useQuery({
    queryKey: getEventStaffUsersQueryKey(eventId),
    queryFn: fetchEventStaffUsers,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
