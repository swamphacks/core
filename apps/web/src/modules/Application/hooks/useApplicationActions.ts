import { api } from "@/lib/ky";
import { useUserQueryKey } from "@/lib/auth/hooks/useUser";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import { myApplicationQueryKey } from "./useMyApplication";

async function performAction(
  method: "post" | "patch",
  endpoint: string,
  failureMessage: string,
) {
  try {
    if (method === "post") {
      await api.post(endpoint);
    } else {
      await api.patch(endpoint);
    }
  } catch (error) {
    let message = failureMessage;
    if (error instanceof HTTPError) {
      const body = (await error.response.json().catch(() => null)) as {
        detail?: string;
        message?: string;
      } | null;
      message = body?.detail || body?.message || failureMessage;
    }
    toast.error(message);
    throw error;
  }
}

export function useApplicationActions() {
  const queryClient = useQueryClient();

  const refreshStatus = async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: myApplicationQueryKey }),
      queryClient.invalidateQueries({ queryKey: useUserQueryKey }),
    ]);
  };

  const confirmAttendance = useMutation({
    mutationFn: () =>
      performAction(
        "patch",
        "application/confirm",
        "Failed to confirm attendance.",
      ),
    onSuccess: async () => {
      toast.success("Attendance confirmed!");
      await refreshStatus();
      window.location.assign("/hacker-portal");
    },
  });

  const withdrawApplication = useMutation({
    mutationFn: () =>
      performAction(
        "patch",
        "application/withdraw",
        "Failed to withdraw attendance.",
      ),
    onSuccess: async () => {
      toast.success("Your attendance has been withdrawn.");
      await refreshStatus();
      window.location.assign("/application");
    },
  });

  const joinWaitlist = useMutation({
    mutationFn: () =>
      performAction(
        "post",
        "application/join-waitlist",
        "Failed to join the waitlist.",
      ),
    onSuccess: async () => {
      toast.success("You are on the waitlist.");
      await refreshStatus();
    },
  });

  const leaveWaitlist = useMutation({
    mutationFn: () =>
      performAction(
        "post",
        "application/leave-waitlist",
        "Failed to leave the waitlist.",
      ),
    onSuccess: async () => {
      toast.success("You have left the waitlist.");
      await refreshStatus();
    },
  });

  return {
    confirmAttendance,
    withdrawApplication,
    joinWaitlist,
    leaveWaitlist,
  };
}
