import { api } from "@/lib/ky";
import { type User } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";
import z from "zod";

export const StaffUserSchema = z.object({
  userId: z.uuid(),
  name: z.string().nullable(),
  email: z.email(),
  image: z.url().nullable(),
  assignedAt: z.iso.datetime(),
  eventRole: z.enum(["staff", "admin"]),
});

export type StaffUser = z.infer<typeof StaffUserSchema>;

export function getEventStaffUsersQueryKey(eventId: string) {
  return ["event", eventId, "staff-users"] as const;
}

export function useEventStaffUsers(eventId: string) {
  async function fetchEventStaffUsers(): Promise<User[]> {
    const result = await api.get<User[]>(`events/${eventId}/staff`).json();
    return result;
  }

  return useQuery({
    queryKey: getEventStaffUsersQueryKey(eventId),
    queryFn: fetchEventStaffUsers,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}
