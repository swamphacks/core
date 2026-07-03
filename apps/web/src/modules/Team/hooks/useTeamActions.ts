import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { myTeamQueryKey } from "./useMyTeam";

type CreateTeamRequest =
  operations["create-team"]["requestBody"]["content"]["application/json"];

type CreateTeamResponse =
  operations["create-team"]["responses"]["200"]["content"]["application/json"];

async function createTeamFn(
  req: CreateTeamRequest,
): Promise<CreateTeamResponse> {
  return api
    .post<CreateTeamResponse>("team", {
      json: req,
    })
    .json();
}

type DeleteTeamRequest =
  operations["delete-team"]["requestBody"]["content"]["application/json"];

async function deleteTeamFn(req: DeleteTeamRequest) {
  return api
    .delete("team", {
      json: req,
    })
    .json();
}

async function joinTeamFn(invitationId: string) {
  return api.post(`team/join/${invitationId}`).json();
}

async function leaveTeamFn(teamId: string) {
  return api.post(`team/${teamId}/leave`).json();
}

type KickMemberRequest =
  operations["kick-member"]["requestBody"]["content"]["application/json"] & {
    teamId: string;
  };

async function kickMemberFn(req: KickMemberRequest) {
  return api
    .post(`team/${req.teamId}/kick`, {
      json: {
        memberId: req.memberId,
      },
    })
    .json();
}

export function useTeamActions() {
  const queryClient = useQueryClient();

  const createTeam = useMutation({
    mutationFn: createTeamFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myTeamQueryKey,
      });
    },
  });

  const deleteTeam = useMutation({
    mutationFn: deleteTeamFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myTeamQueryKey,
      });
    },
  });

  const joinTeam = useMutation({
    mutationFn: joinTeamFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myTeamQueryKey,
      });
    },
  });

  const leaveTeam = useMutation({
    mutationFn: leaveTeamFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myTeamQueryKey,
      });
    },
  });

  const kickMember = useMutation({
    mutationFn: kickMemberFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myTeamQueryKey,
      });
    },
  });

  return { createTeam, deleteTeam, joinTeam, leaveTeam, kickMember };
}
