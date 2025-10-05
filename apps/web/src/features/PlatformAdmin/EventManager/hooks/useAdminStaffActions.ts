import { useMutation, useQueryClient } from "@tanstack/react-query";
import { getEventStaffUsersQueryKey } from "./useEventStaffUsers";
import z from "zod";
import { api } from "@/lib/ky";

export const assignStaffRoleSchema = z.object({
  email: z.email("Must enter an email."),
  role: z.enum(["staff", "admin"], "Must be a valid role"),
});

export type AssignStaffRole = z.infer<typeof assignStaffRoleSchema>;

export function useAdminStaffActions(eventId: string) {
  const queryClient = useQueryClient();

  async function addEventStaff(data: AssignStaffRole) {
    await api.post(`events/${eventId}/roles`, { json: data });
  }

  const add = useMutation({
    mutationFn: addEventStaff,
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: getEventStaffUsersQueryKey(eventId),
      });
    },
    onError: (error) => {
      console.log(error);
    },
  });

  return { add };
}
