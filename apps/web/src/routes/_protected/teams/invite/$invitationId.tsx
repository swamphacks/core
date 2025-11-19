import { createFileRoute, redirect, useRouter } from "@tanstack/react-router";
import { useQuery, useMutation } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import { HTTPError } from "ky";
import { Heading, Text } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { PageLoading } from "@/components/PageLoading";
import { useTheme } from "@/components/ThemeProvider";
import { toast } from "react-toastify";
import ArrowIcon from "~icons/tabler/arrow-right";
import { auth } from "@/lib/authClient";
import { getUserEventRole } from "@/features/Event/api/getUserEventRole";
import { fetchApplication } from "@/features/Application/hooks/useApplication";

interface TeamMember {
  user_id: string;
  name: string;
  email: string | null;
  image: string | null;
  joined_at: string | null;
}

interface InvitationDetails {
  id: string;
  team_name: string;
  inviter_name: string;
  event_name: string;
  event_id: string;
  invited_email: string;
  status: string;
  expires_at: string | null;
  created_at: string;
  team_members: TeamMember[];
}

async function fetchInvitation(invitationId: string): Promise<InvitationDetails> {
  try {
    const response = await api
      .get<InvitationDetails>(`teams/invite/${invitationId}`)
      .json();
    return response;
  } catch (error) {
    if (error instanceof HTTPError && error.response.status === 404) {
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

async function linkUserToInvitation(invitationId: string): Promise<{ expired: boolean }> {
  const response = await api.post<{ expired: boolean }>(`teams/invite/${invitationId}/claim`).json();
  return response;
}

interface BeforeLoadData {
  valid: boolean;
  expired?: boolean;
  eventId?: string;
}

export const Route = createFileRoute("/_protected/teams/invite/$invitationId")({
  component: RouteComponent,
  beforeLoad: async ({ context, location, params }): Promise<BeforeLoadData> => {
    const { user } = context;
    if (!user) {
      throw redirect({
        to: "/",
        search: { redirect: location.pathname },
      });
    }

    const { invitationId } = params;

    let expired = false;
    try {
      const linkResult = await linkUserToInvitation(invitationId);
      expired = linkResult.expired;
    } catch (error) {
      if (error instanceof HTTPError) {
        if (error.response.status === 404) {
          expired = true;
        } else {
          console.error("Error linking user to invitation:", error);
        }
      } else {
        console.error("Error linking user to invitation:", error);
      }
    }

    if (expired) {
      return {
        valid: true,
        expired: true,
      };
    }

    let eventId: string | undefined;
    try {
      const invitationDetails = await fetchInvitation(invitationId);
      eventId = invitationDetails.event_id;
    } catch (error) {
      if (error instanceof HTTPError && error.response.status === 404) {
        return {
          valid: true,
          expired: true,
        };
      }
      throw error;
    }

    if (!eventId) {
      return {
        valid: true,
        expired: true,
      };
    }

    let hasEventAccess = false;

    try {
      const eventRole = await getUserEventRole(eventId);
      if (eventRole?.role === "attendee" || eventRole?.role === "applicant") {
        hasEventAccess = true;
      }

      if (!hasEventAccess) {
        try {
          const application = await fetchApplication(eventId);
          if (application?.status === "submitted" || application?.status === "under_review" || 
              application?.status === "accepted" || application?.status === "waitlisted") {
            hasEventAccess = true;
          }
        } catch (error) {
          if (error instanceof HTTPError && error.response.status !== 404) {
            console.warn("Error fetching application for event access check:", error);
          }
        }
      }
    } catch (error) {
      if (error instanceof HTTPError && error.response.status !== 404) {
        console.warn("Error checking event role for access:", error);
      }
    }

    if (!hasEventAccess) {
      toast.warning("You need to apply to this event to join a team.", {
        position: "bottom-right",
      });
      throw redirect({
        to: "/events/$eventId/application",
        params: { eventId },
        search: { redirect: location.pathname } as { redirect: string },
      });
    }

    return {
      valid: true,
      expired: false,
      eventId,
    };
  }
});

function ExpiredInvitationComponent() {
  const router = useRouter();
  const { theme } = useTheme();

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-background">
      <div className="max-w-2xl w-full flex flex-col items-center">
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

        <Heading className="text-3xl font-medium mb-8 text-text-main">
          Invitation Not Available
        </Heading>

        <div className="mb-6 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md max-w-md">
          <Text className="text-yellow-800 dark:text-yellow-200 text-center">
            This invitation has expired, been already accepted, or doesn't exist.
          </Text>
        </div>

        <Button
          onPress={() => router.navigate({ to: "/portal" })}
          variant="primary"
          className="w-full max-w-md"
        >
          Go to Portal
        </Button>
      </div>
    </div>
  );
}

function RouteComponent() {
  const { invitationId } = Route.useParams();
  const router = useRouter();
  const beforeLoadData = Route.useRouteContext();
  const { theme } = useTheme();

  const expired = beforeLoadData?.expired === true;

  const invitationQuery = useQuery({
    queryKey: ["invitation", invitationId],
    queryFn: () => fetchInvitation(invitationId),
    enabled: !expired,
  });

  const userQuery = auth.useUser();
  const currentUser = userQuery.data?.user;

  const acceptMutation = useMutation({
    mutationFn: () => acceptInvitation(invitationId),
    onSuccess: () => {
      toast.success("Invitation accepted! You've been added to the team.", {
        position: "bottom-right",
      });
      router.navigate({ to: "/portal" });
    },
    onError: async (error: unknown) => {
      try {
        if (error instanceof HTTPError) {
          // user not authenticated redirect to login
          if (error.response.status === 401) {
            const invitationUrl = `/teams/invite/${invitationId}`;
            router.navigate({
              to: "/",
              search: { redirect: invitationUrl },
            });
            toast.info("Please log in to accept the invitation");
            return;
          }

          const data = await error.response.json();
          
          if (error.response.status === 403 && data.error === "application_required") {
            const eventId = invitationQuery.data?.event_id;
            if (eventId) {
              toast.error(
                <div className="flex flex-col gap-2">
                  <Text className="text-sm">{data.message || "Please complete your application first."}</Text>
                  <Button
                    variant="primary"
                    size="sm"
                    onPress={() => {
                      const inviteLink = `/teams/invite/${invitationId}`;
                      router.navigate({
                        to: "/events/$eventId/application",
                        params: { eventId },
                        search: { redirect: inviteLink } as { redirect: string },
                      });
                    }}
                    className="self-start"
                  >
                    Go to Application
                  </Button>
                </div>,
                {
                  autoClose: 10000,
                }
              );
              return;
            }
          }

          toast.error(data.message || "Failed to accept invitation");
        } else if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error("Failed to accept invitation");
        }
      } catch {
        toast.error("Failed to accept invitation");
      }
    },
  });

  const rejectMutation = useMutation({
    mutationFn: () => rejectInvitation(invitationId),
    onSuccess: () => {
      toast.success("Invitation rejected.");
      router.navigate({ to: "/portal" });
    },
    onError: async (error: unknown) => {
      try {
        if (error instanceof HTTPError) {
          if (error.response.status === 401) {
            const invitationUrl = `/teams/invite/${invitationId}`;
            router.navigate({
              to: "/",
              search: { redirect: invitationUrl },
            });
            toast.info("Please log in to reject the invitation");
            return;
          }

          const data = await error.response.json();
          toast.error(data.message || "Failed to reject invitation");
        } else if (error instanceof Error) {
          toast.error(error.message);
        } else {
          toast.error("Failed to reject invitation");
        }
      } catch {
        toast.error("Failed to reject invitation");
      }
    },
  });

  if (expired) {
    return <ExpiredInvitationComponent />;
  }

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
  const isExpired =
    invitation.expires_at &&
    new Date(invitation.expires_at) < new Date();
  const isAccepted = invitation.status === "ACCEPTED";
  const isRejected = invitation.status === "REJECTED";

  return (
    <div className="min-h-screen flex items-center justify-center p-6 bg-background">
      <div className="max-w-2xl w-full flex flex-col items-center">
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

        <Heading className="text-3xl font-medium mb-8 text-text-main">
          SwampHacks Team Invitation
        </Heading>

        <div className="flex items-center justify-center mb-8 relative w-full max-w-md">

          <div className="w-40 h-40 rounded-full overflow-hidden bg-[#929292] dark:bg-gray-600 flex items-center justify-center border-4 border-white dark:border-gray-700 shadow-lg">
            {currentUser?.image ? (
              <img
                src={currentUser.image}
                alt={currentUser.name || "You"}
                className="w-full h-full object-cover"
                onError={(e) => {
                  const target = e.target as HTMLImageElement;
                  target.style.display = 'none';
                  const parent = target.parentElement;
                  if (parent) {
                    const initial = document.createElement('span');
                    initial.className = 'text-white text-2xl font-medium';
                    initial.textContent = (currentUser.name || 'U').charAt(0).toUpperCase();
                    parent.appendChild(initial);
                  }
                }}
              />
            ) : currentUser?.name ? (
              <span className="text-white text-2xl font-medium">
                {currentUser.name.charAt(0).toUpperCase()}
              </span>
            ) : (
              <span className="text-white text-2xl font-medium">?</span>
            )}
          </div>

          <div className="mx-8 flex items-center">
            <ArrowIcon
              className="h-15 w-15 text-text-main"
            />
          </div>

          <div className="grid grid-cols-2 gap-2">
            {invitation.team_members && invitation.team_members.length > 0 ? (
              invitation.team_members.slice(0, 4).map((member) => (
                <div
                  key={member.user_id}
                  className="w-16 h-16 rounded-full overflow-hidden bg-[#929292] dark:bg-gray-600 flex items-center justify-center border-2 border-white dark:border-gray-700 shadow-sm"
                  title={member.name}
                >
                  {member.image ? (
                    <img
                      src={member.image}
                      alt={member.name}
                      className="w-full h-full object-cover"
                      onError={(e) => {
                        const target = e.target as HTMLImageElement;
                        target.style.display = 'none';
                        const parent = target.parentElement;
                        if (parent) {
                          const initial = document.createElement('span');
                          initial.className = 'text-white text-xs font-medium';
                          initial.textContent = member.name.charAt(0).toUpperCase();
                          parent.appendChild(initial);
                        }
                      }}
                    />
                  ) : (
                    <span className="text-white text-xs font-medium">
                      {member.name.charAt(0).toUpperCase()}
                    </span>
                  )}
                </div>
              ))
            ) : (
              (<>
                <div className="w-16 h-16 bg-[#929292] dark:bg-gray-600 rounded-full"></div>
                <div className="w-16 h-16 bg-[#929292] dark:bg-gray-600 rounded-full"></div>
                <div className="w-16 h-16 bg-[#929292] dark:bg-gray-600 rounded-full"></div>
                <div className="w-16 h-16 bg-[#929292] dark:bg-gray-600 rounded-full"></div>
              </>)
            )}
          </div>
        </div>

        <Text className="text-lg mb-8 text-text-main text-center">
          <strong>{invitation.inviter_name}</strong> has invited you to join their
          team for <strong>{invitation.event_name}</strong>
        </Text>

        {isExpired && (
          <div className="mb-6 p-4 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-md">
            <Text className="text-yellow-800 dark:text-yellow-200">
              This invitation has expired.
            </Text>
          </div>
        )}

        {isAccepted && (
          <div className="mb-6 p-4 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-md">
            <Text className="text-green-800 dark:text-green-200">
              This invitation has already been accepted.
            </Text>
          </div>
        )}

        {isRejected && (
          <div className="mb-6 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md">
            <Text className="text-red-800 dark:text-red-200">
              This invitation has been rejected.
            </Text>
          </div>
        )}

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

        <div className="mt-8 py-4 px-15 bg-surface border-1 border-border rounded-md max-w-lg text-center">
          <Text className="text-sm text-text-secondary">
            You may always leave and join teams at your own discretion anytime
            during the event.
          </Text>
        </div>
      </div>
    </div>
  )
}
