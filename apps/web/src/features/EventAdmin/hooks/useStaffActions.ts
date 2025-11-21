import {
  getEventStaffUsersQueryKey,
  type StaffUsers,
} from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { api } from "@/lib/ky";
import type { components } from "@/lib/openapi/schema";
import { useMutation, useQueryClient } from "@tanstack/react-query";

type AddRoleFields = {
  assignments: (Omit<
    components["schemas"]["handlers.AssignRoleFields"],
    "email"
  > & { email?: string })[];
};

async function addEventStaff(eventId: string, userIds: string[]) {
  if (userIds.length === 0) return;

  const body: AddRoleFields = {
    assignments: userIds.map((userId) => ({
      user_id: userId,
      role: "staff",
    })),
  };

  await api.post(`events/${eventId}/roles/batch`, {
    json: body,
  });
}

async function deleteEventStaff(eventId: string, userId: string) {
  await api.delete(`events/${eventId}/roles/${userId}`);

  return userId;
}

export function useStaffActions(eventId: string) {
  const queryClient = useQueryClient();

  const addMany = useMutation({
    mutationFn: (userIds: string[]) => addEventStaff(eventId, userIds),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: getEventStaffUsersQueryKey(eventId),
      });
    },
  });

  const remove = useMutation({
    mutationFn: (userId: string) => deleteEventStaff(eventId, userId),
    onSuccess: (userId) => {
      queryClient.setQueryData<StaffUsers>(
        getEventStaffUsersQueryKey(eventId),
        (old) => {
          if (!old) return old;

          return old.filter((user) => user.id !== userId);
        },
      );
    },
  });

  return {
    addMany,
    remove,
  };
}
