import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getEventStaffUsersQueryKey,
  type StaffUser,
} from "./useEventStaffUsers";
import z from "zod";

export const addStaffSchema = z.object({
  email: z.email("Must enter an email."),
  role: z.enum(["STAFF", "ADMIN"], "Must be a valid role"),
});

export type AddStaff = z.infer<typeof addStaffSchema>;

async function addEventStaff(data: AddStaff): Promise<StaffUser> {
  // Needs to make POST to backend
  return {
    userId: "6c62875a-a52c-4e98-8b02-5ead9ec0e193",
    name: "Ben Gerald",
    email: data.email,
    eventRole: data.role,
    image: null,
    assignedAt: "2025-08-03T18:45:30.123Z",
  };
}

export function useAdminStaffActions(eventId: string) {
  const queryClient = useQueryClient();

  const add = useMutation({
    mutationFn: addEventStaff,
    onSuccess: (addedStaff) => {
      queryClient.setQueryData<StaffUser[]>(
        getEventStaffUsersQueryKey(eventId),
        (old) => {
          if (!old) return [addedStaff];

          return [...old, addedStaff];
        },
      );
    },
  });

  return { add };
}
