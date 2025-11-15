import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import z from "zod";

export const newTeamSchema = z.object({
  name: z
    .string()
    .min(2, "Must be more than 2 characters.")
    .max(24, "Must be less than 24 characters."),
});
export type NewTeam = z.infer<typeof newTeamSchema>;

async function leaveTeam(teamId: string) {
  await api.delete(`teams/${teamId}/members/me`);
}

async function createTeam(eventId: string, data: NewTeam) {
  try {
    await api.post(`events/${eventId}/teams`, {
      json: data,
    });
  } catch (err) {
    if (err instanceof HTTPError && err.response.status === 409) {
      toast.error("You are already part of a team. Leave it to make a team.");
    }
    throw err; // rethrow other errors so Tanstack Query can handle them
  }
}

export function useTeamActions(eventId: string) {
  const queryClient = useQueryClient();

  const leave = useMutation({
    mutationFn: (teamId: string) => leaveTeam(teamId),
    onSuccess: () => {
      queryClient.setQueryData(["myTeam", eventId], null);
    },
  });

  const create = useMutation({
    mutationFn: (data: NewTeam) => createTeam(eventId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["myTeam", eventId],
      });
    },
  });

  return { leave, create };
}
