import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import { getMyPendingJoinRequestsQueryKey } from "./useMyPendingJoinRequests";
import type { ErrorResponse } from "@/lib/openapi/types";
import { getTeamPendingJoinRequestsQueryKey } from "./useTeamPendingJoinRequests";

async function createJoinRequest(teamId: string, eventId: string) {
  try {
    await api.post(`events/${eventId}/teams/${teamId}/join`, {
      json: {
        message: null, // No message for now, TODO: allow users to add a message
      },
    });
  } catch (err) {
    if (err instanceof HTTPError && err.response.status === 409) {
      toast.error(
        "You currently have a conflicting join request or team membership.",
      );
    }

    throw err;
  }

  return eventId;
}

async function acceptJoinRequest(requestId: string, teamId: string) {
  try {
    await api.post(`teams/join/${requestId}/accept`);
  } catch (err) {
    if (err instanceof HTTPError) {
      // Parse payload to get error message
      const errorBody = await err.response.json<ErrorResponse>();
      toast.error(errorBody.message || "An error occurred.");
    } else {
      toast.error("An error occurred.");
    }

    throw err;
  }

  // Return teamId to use in onSuccess callback
  return teamId;
}

async function rejectJoinRequest(requestId: string, teamId: string) {
  try {
    await api.post(`teams/join/${requestId}/reject`);
  } catch (err) {
    if (err instanceof HTTPError) {
      // Parse payload to get error message
      const errorBody = await err.response.json<ErrorResponse>();
      toast.error(errorBody.message || "An error occurred.");
    } else {
      toast.error("An error occurred.");
    }

    throw err;
  }

  // Return teamId to use in onSuccess callback
  return teamId;
}

export function useJoinRequestActions() {
  const queryClient = useQueryClient();

  const create = useMutation({
    mutationFn: ({ teamId, eventId }: { teamId: string; eventId: string }) =>
      createJoinRequest(teamId, eventId),
    onSuccess: (eventId) => {
      queryClient.invalidateQueries({
        queryKey: getMyPendingJoinRequestsQueryKey(eventId),
      });
    },
  });

  const acceptRequest = useMutation({
    mutationFn: ({
      requestId,
      teamId,
    }: {
      requestId: string;
      teamId: string;
    }) => acceptJoinRequest(requestId, teamId),
    onSuccess: (teamId) => {
      queryClient.invalidateQueries({
        queryKey: getTeamPendingJoinRequestsQueryKey(teamId),
      });
    },
  });

  const rejectRequest = useMutation({
    mutationFn: ({
      requestId,
      teamId,
    }: {
      requestId: string;
      teamId: string;
    }) => rejectJoinRequest(requestId, teamId),
    onSuccess: (teamId) => {
      queryClient.invalidateQueries({
        queryKey: getTeamPendingJoinRequestsQueryKey(teamId),
      });
    },
  });

  return { create, acceptRequest, rejectRequest };
}
