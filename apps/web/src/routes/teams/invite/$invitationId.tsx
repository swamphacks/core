import { createFileRoute, useRouter } from "@tanstack/react-router";
import { useQuery, useMutation } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import { HTTPError } from "ky";
import { Heading, Text } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { PageLoading } from "@/components/PageLoading";
import { useTheme } from "@/components/ThemeProvider";
import { toast } from "react-toastify";
import ArrowIcon from "~icons/tabler/arrow-right";

interface InvitationDetails {
  id: string;
  team_name: string;
  inviter_name: string;
  event_name: string;
  invited_email: string;
  status: string;
  expires_at: string | null;
  created_at: string;
}

async function fetchInvitation(invitationId: string): Promise<InvitationDetails> {
  try {
    const response = await api
      .get<InvitationDetails>(`teams/invite/${invitationId}`)
      .json();
    return response;
  } catch (error) {
    // Re-throw to let React Query handle it
    if (error instanceof HTTPError && error.response.status === 404) {
      // Still throw so React Query knows it's an error, but we'll handle it in the component
      throw error;
    }
    throw error;
  }
}

async function acceptInvitation(invitationId: string): Promise<void> {
  await api.post(`teams/invite/${invitationId}/accept`).json();
}

async function rejectInvitation(invitationId: string): Promise<void> {
  await api.post(`teams/invite/${invitationId}/reject`).json();
}

export const Route = createFileRoute("/teams/invite/$invitationId")({
  component: RouteComponent,
});

function RouteComponent() {
  const { invitationId } = Route.useParams();
  const router = useRouter();

  const invitationQuery = useQuery({
    queryKey: ["invitation", invitationId],
    queryFn: () => fetchInvitation(invitationId),
  });

  const acceptMutation = useMutation({
    mutationFn: () => acceptInvitation(invitationId),
    onSuccess: () => {
      toast.success("Invitation accepted! You've been added to the team.");
      router.navigate({ to: "/portal" });
    },
    onError: async (error: unknown) => {
      let message = "Failed to accept invitation";
      try {
        if (error instanceof HTTPError) {
          const data = await error.response.json();
          message = data.message || message;
        } else if (error instanceof Error) {
          message = error.message;
        }
      } catch {
        // Use default message
      }
      toast.error(message);
    },
  });

  const rejectMutation = useMutation({
    mutationFn: () => rejectInvitation(invitationId),
    onSuccess: () => {
      toast.success("Invitation rejected.");
      router.navigate({ to: "/portal" });
    },
    onError: async (error: unknown) => {
      let message = "Failed to reject invitation";
      try {
        if (error instanceof HTTPError) {
          const data = await error.response.json();
          message = data.message || message;
        } else if (error instanceof Error) {
          message = error.message;
        }
      } catch {
        // Use default message
      }
      toast.error(message);
    },
  });

  if (invitationQuery.isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <PageLoading />
      </div>
    );
  }

  if (invitationQuery.isError || !invitationQuery.data) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="max-w-md mx-auto p-6 text-center">
          <Heading className="text-2xl mb-4">Invitation Not Found</Heading>
          <Text className="text-text-secondary mb-6">
            This invitation may have expired, been already accepted, or doesn't
            exist.
          </Text>
          <Button
            onPress={() => router.navigate({ to: "/portal" })}
            variant="primary"
          >
            Go to Portal
          </Button>
        </div>
      </div>
    );
  }

  const invitation = invitationQuery.data;
  const { theme } = useTheme();
  const isExpired =
    invitation.expires_at &&
    new Date(invitation.expires_at) < new Date();
  const isAccepted = invitation.status === "ACCEPTED";
  const isRejected = invitation.status === "REJECTED";

  // Get first letter of inviter's name for the avatar
  const inviterInitial = invitation.inviter_name.charAt(0).toUpperCase();

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-white">
      <div className="max-w-2xl w-full flex flex-col items-center">
        {/* Logo */}
        <div className="mb-8 flex justify-center">
          <img
            src={
              theme === "dark"
                ? "/assets/SwampHacks_Logo_Light.png"
                : "/assets/SwampHacks_Logo_Dark.svg"
            }
            alt="SwampHacks Logo"
            className="h-16"
          />
        </div>

        {/* Title */}
        <Heading className="text-3xl font-medium mb-8 text-black">
          SwampHacks Team Invitation
        </Heading>

        {/* Visual Representation */}
        <div className="flex items-center justify-center mb-8 relative w-full max-w-md">

          {/* Large gray circle (individual) */}
          <div className="w-40 h-40 bg-[#929292] rounded-full"></div>

          {/* Arrow */}
          <div className="mx-8 flex items-center">
            <ArrowIcon
              className="h-15 w-15"
            />
          </div>

          {/* Four smaller circles (team) */}
          <div className="grid grid-cols-2 gap-2">
            <div className="w-15 h-15 bg-[#929292] rounded-full"></div>
            <div className="w-15 h-15 bg-[#929292] rounded-full"></div>
            <div className="w-15 h-15 bg-[#929292] rounded-full"></div>
            <div className="w-15 h-15 bg-[#929292] rounded-full"></div>
          </div>
        </div>

        {/* Invitation Text */}
        <Text className="text-lg mb-8 text-black text-center">
          <strong>{invitation.inviter_name}</strong> has invited you to join their
          team for <strong>{invitation.event_name}</strong>
        </Text>

        {/* Status Messages */}
        {isExpired && (
          <div className="mb-6 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <Text className="text-yellow-800">
              This invitation has expired.
            </Text>
          </div>
        )}

        {isAccepted && (
          <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-md">
            <Text className="text-green-800">
              This invitation has already been accepted.
            </Text>
          </div>
        )}

        {isRejected && (
          <div className="mb-6 p-4 bg-red-50 border border-red-200 rounded-md">
            <Text className="text-red-800">
              This invitation has been rejected.
            </Text>
          </div>
        )}

        {/* Action Buttons */}
        {!isExpired && !isAccepted && !isRejected && (
          <div className="flex flex-row gap-4 mb-8 w-full max-w-md">
            <Button
              onPress={() => acceptMutation.mutate()}
              isDisabled={acceptMutation.isPending || rejectMutation.isPending}
              className="flex-1 !bg-[#B9F8CF] hover:!bg-[#a2eabb] text-[#00A63E] border-1 border-[#7df7aa] font-light py-3 px-6"
              style={{
                backgroundColor: "#86efac",
              }}
            >
              {acceptMutation.isPending ? "Accepting..." : "Accept"}
            </Button>
            <Button
              onPress={() => rejectMutation.mutate()}
              isDisabled={acceptMutation.isPending || rejectMutation.isPending}
              className="flex-1 !bg-[#FFC9C9] hover:!bg-[#fab1b1] text-[#E7000B] border-1 border-[#f79e9e] font-light py-3 px-6"
              style={{
                backgroundColor: "#fca5a5",
              }}
            >
              {rejectMutation.isPending ? "Rejecting..." : "Deny"}
            </Button>
          </div>
        )}

        {(isExpired || isAccepted || isRejected) && (
          <Button
            onPress={() => router.navigate({ to: "/portal" })}
            className="w-full max-w-md !bg-green-400 hover:!bg-green-500 text-white font-semibold py-3 px-6 mb-8"
            style={{
              backgroundColor: "#86efac",
            }}
          >
            Go to Portal
          </Button>
        )}

        {/* Disclaimer */}
        <div className="mt-8 py-4 px-15 bg-gray-100 border-1 border-[#D4D4D8] rounded-md max-w-lg text-center">
          <Text className="text-sm text-gray-700">
            You may always leave and join teams at your own discretion anytime
            during the event.
          </Text>
        </div>
      </div>
    </div>
  );
}
