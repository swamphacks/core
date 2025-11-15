import { Heading } from "react-aria-components";
import TeamJoinRequestCard from "./TeamJoinRequestCard";
import TablerGitPullRequest from "~icons/tabler/git-pull-request";
import { useTeamPendingJoinRequests } from "../hooks/useTeamPendingJoinRequests";
import { useJoinRequestActions } from "../hooks/useJoinRequestActions";
import { toast } from "react-toastify";

interface Props {
  teamId: string;
}

export default function TeamJoinRequestSection({ teamId }: Props) {
  const joinRequests = useTeamPendingJoinRequests(teamId);
  const { acceptRequest, rejectRequest } = useJoinRequestActions();

  const handleAccept = (requestId: string) => {
    acceptRequest.mutate(
      { requestId, teamId },
      {
        onSuccess: () => {
          toast.success("Join request accepted.");
        },
      },
    );
  };

  const handleReject = (requestId: string) => {
    rejectRequest.mutate(
      { requestId, teamId },
      {
        onSuccess: () => {
          toast.success("Join request rejected.");
        },
      },
    );
  };

  return (
    <section>
      <Heading className="text-2xl text-text-main flex flex-row gap-2 mb-4 items-center">
        <TablerGitPullRequest className="w-5 h-5" />
        Team Join Requests
      </Heading>

      <div className="flex flex-col gap-4">
        {joinRequests.data && joinRequests.data.length > 0 ? (
          joinRequests.data.map((request) => (
            <TeamJoinRequestCard
              onAccept={() => handleAccept(request.id)}
              onReject={() => handleReject(request.id)}
              requesterName={request.user_name}
            />
          ))
        ) : (
          <p className="text-text-secondary">No pending join requests.</p>
        )}
      </div>
    </section>
  );
}
