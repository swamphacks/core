import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";

async function leaveTeam(teamId: string) {
  await api.delete(`teams/${teamId}/members/me`);
}

export function useTeamActions(eventId: string) {
  const queryClient = useQueryClient();

  const leave = useMutation({
    mutationFn: (teamId: string) => leaveTeam(teamId),
    onSuccess: () => {
      queryClient.setQueryData(["myTeam", eventId], null);
    },
  });

  return { leave };
}
