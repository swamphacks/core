import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import { myApplicationQueryKey } from "./useMyApplication";
import type { ErrorResponse } from "@/lib/auth/types";

async function confirmAttendanceFn() {
  try {
    await api.patch("application/confirm");
  } catch (err) {
    if (err instanceof HTTPError) {
      const errorBody = await err.response.json<ErrorResponse>();
      toast.error(errorBody.message || "Failed to confirm attendance.");
    } else {
      toast.error("An error occurred while confirming attendance.");
    }

    throw err;
  }
}

async function withdrawApplicationFn() {
  try {
    await api.patch("application/withdraw");
  } catch (err) {
    if (err instanceof HTTPError) {
      const errorBody = await err.response.json<ErrorResponse>();
      toast.error(errorBody.message || "Failed to withdraw application.");
    } else {
      toast.error("An error occurred while withdrawing application.");
    }

    throw err;
  }
}

export function useApplicationActions() {
  const queryClient = useQueryClient();

  const confirmAttendance = useMutation({
    mutationFn: confirmAttendanceFn,
    onSuccess: () => {
      toast.success("Attendance confirmed!");
      queryClient.invalidateQueries({
        queryKey: myApplicationQueryKey,
      });
    },
  });

  const withdrawApplication = useMutation({
    mutationFn: withdrawApplicationFn,
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: myApplicationQueryKey,
      });
    },
  });

  return { confirmAttendance, withdrawApplication };
}
