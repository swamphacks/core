import { AvatarStack } from "@/components/ui/AvatarStack";
import type { TeamWithMembers } from "../hooks/useMyTeam";
import TablerUsers from "~icons/tabler/users";
import TablerCheck from "~icons/tabler/check";
import { Button } from "@/components/ui/Button";
import { useJoinRequestActions } from "../hooks/useJoinRequestActions";
import { toast } from "react-toastify";

interface Props {
  team: TeamWithMembers;
  isCurrentTeam?: boolean;
  alreadyRequested?: boolean;
}

export default function TeamCard({
  team,
  isCurrentTeam = false,
  alreadyRequested = false,
}: Props) {
  const { create } = useJoinRequestActions();
  const handleJoinRequest = async () => {
    await create.mutateAsync(
      {
        teamId: team.id,
        eventId: team.event_id,
      },
      {
        onSuccess: () => {
          toast.success("Join request sent!");
        },
        onError: () => {
          toast.error("Failed to send join request. Try again later.");
        },
      },
    );
  };

  return (
    <div className="flex flex-col bg-surface md:max-w-86 w-full px-6 py-4 rounded-sm gap-2">
      <h2 className="text-xl font-medium">{team.name}</h2>
      <div className="flex flex-row text-text-secondary items-center gap-4">
        <AvatarStack
          avatars={team.members.map((member) => ({
            src: member?.image,
            fallback: member?.name,
          }))}
          max={2}
          size="sm"
        />
        <div className="flex flex-row justify-center items-center gap-2">
          <TablerUsers className="h-6 w-6" />
          <p className="text-lg">{team.members.length} / 4 Members</p>
        </div>
      </div>

      {isCurrentTeam ? (
        <Button isDisabled className="mt-4">
          Your Current Team
        </Button>
      ) : alreadyRequested ? (
        <Button
          variant="secondary"
          className="mt-4 bg-green-500 hover:bg-green-500 pressed:bg-green-500 hover:cursor-default text-white rounded-xs flex items-center gap-2"
        >
          <TablerCheck className="h-5 w-5" />
          Requested
        </Button>
      ) : (
        <Button
          onClick={handleJoinRequest}
          variant="secondary"
          className="mt-4 rounded-xs"
        >
          Request to Join
        </Button>
      )}
    </div>
  );
}
