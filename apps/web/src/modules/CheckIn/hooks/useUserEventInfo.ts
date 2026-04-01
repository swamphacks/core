import { api } from "@/lib/ky";
import { EventRoleSchema, PlatformRoleSchema } from "@/lib/openapi/zodSchemas";
import { useQuery } from "@tanstack/react-query";
import z from "zod";

const userEventInfoSchema = z.object({
  user_id: z.uuid(),
  name: z.string(),
  email: z.email(),
  image: z.url().nullable(),
  platform_role: PlatformRoleSchema,
  event_role: EventRoleSchema.nullable(),
  checked_in_at: z.coerce.date().nullable(),
});

export type UserEventInfo = z.infer<typeof userEventInfoSchema>;

const fetchUserEventInfo = async (userId: string | null, eventId: string) => {
  if (!userId) {
    return null;
  }
  const result = await api.get(`events/${eventId}/users/${userId}`).json();

  return userEventInfoSchema.parse(result);
};

/**
 *
 * @param userId - ID of the user to fetch info for
 * @param eventId - ID of the event
 * @returns - Object containing user event info. See UserEventInfo type.
 */
export const useUserEventInfo = (eventId: string, userId: string | null) => {
  return useQuery({
    queryKey: ["userEventInfo", userId],
    queryFn: () => fetchUserEventInfo(userId, eventId),
    enabled: !!userId,
  });
};
