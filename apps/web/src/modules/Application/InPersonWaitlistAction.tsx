import { Button } from "@/components/ui/Button";
import { _useUser } from "@/lib/auth/hooks/useUser";
import { api } from "@/lib/ky";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";
import { toast } from "react-toastify";

interface Props {
  userId: string;
  status: string;
}

export default function InPersonWaitlistAction({ userId, status }: Props) {
  const currentUser = _useUser();
  const queryClient = useQueryClient();

  const recordArrival = useMutation({
    mutationFn: async () => {
      try {
        await api.post("application/waitlist/in-person", {
          json: { userId },
        });
      } catch (error) {
        let message = "Unable to record in-person arrival.";
        if (error instanceof HTTPError) {
          const body = (await error.response.json().catch(() => null)) as {
            detail?: string;
            message?: string;
          } | null;
          message = body?.detail || body?.message || message;
        }
        toast.error(message);
        throw error;
      }
    },
    onSuccess: async () => {
      toast.success("In-person arrival recorded.");
      await queryClient.invalidateQueries();
    },
  });

  if (
    currentUser.data?.user?.role !== "admin" ||
    !["rejected", "waitlisted"].includes(status)
  ) {
    return null;
  }

  return (
    <section className="space-y-3 rounded-md border p-4">
      <h3 className="font-semibold">In-person waitlist</h3>
      <p className="text-sm">
        Record arrival only after this hacker is physically present. This gives
        them priority in the event-day waitlist without immediately accepting
        them.
      </p>
      <Button
        onClick={() => {
          if (
            window.confirm(
              "Confirm this hacker is physically present and wants to join the in-person waitlist?",
            )
          ) {
            recordArrival.mutate();
          }
        }}
        isDisabled={recordArrival.isPending}
      >
        {recordArrival.isPending ? "Recording..." : "Record In-Person Arrival"}
      </Button>
      <p className="text-sm">
        The server allows this action only while the in-person waitlist is open.
        Recording arrival again preserves the original arrival time.
      </p>
    </section>
  );
}
