import { AvatarStack } from "@/components/ui/AvatarStack";
import TablerUsers from "~icons/tabler/users";
import type { TeamWithMembers } from "../hooks/useMyTeam";
import { useTeamActions } from "../hooks/useTeamActions";
import { Button } from "@/components/ui/Button";
import { toast } from "react-toastify";
import TablerDoorExit from "~icons/tabler/door-exit";
import TablerCircleX from "~icons/tabler/circle-x";

interface Props {
  eventId: string;
  userId: string;
  team: TeamWithMembers;
}

export default function MyTeamCard({ eventId, userId, team }: Props) {
  const { leave, kickTeamMember } = useTeamActions(eventId);

  const handleLeaveTeam = () => {
    leave.mutate(team.id, {
      onSuccess: () => {
        toast.success("Left the team successfully.");
      },
      onError: () => {
        toast.error("Something went wrong, try again.");
      },
    });
  };

  const handleKickMember = (memberId: string) => {
    kickTeamMember.mutate(
      { teamId: team.id, memberId },
      {
        onSuccess: () => {
          toast.success("Member removed successfully.");
        },
        onError: () => {
          toast.error("Failed to remove member. Try again later.");
        },
      },
    );
  };

  return (
    <div className="border border-input-border rounded-md px-4 py-3 w-full md:w-116 flex flex-col gap-4">
      <div className="w-full flex flex-row items-center justify-between">
        <h3 className="text-text-main text-2xl">{team.name}</h3>

        <Button
          onClick={handleLeaveTeam}
          variant="icon"
          className="aspect-square p-2.5"
        >
          <TablerDoorExit className="w-5 h-5 text-red-600" />
        </Button>
      </div>
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

      <div className="flex flex-col">
        <p>Members:</p>

        <ul className="list-disc list-outside">
          {team.members.map((member) => (
            <li
              key={member.user_id}
              className="text-text-secondary flex items-center justify-between"
            >
              <span>{member.name}</span>

              {/* Only show kick button if the viewer is the owner and the member is not the owner */}
              {member.user_id !== team.owner_id && userId == team.owner_id && (
                <button
                  type="button"
                  onClick={() => handleKickMember(member.user_id)}
                  className="ml-2 p-1 hover:text-red-500 cursor-pointer transition-all duration-150"
                  aria-label={`Remove ${member.name}`}
                >
                  <TablerCircleX className="w-4 h-4" />
                </button>
              )}

              {/* Indicate the owner */}
              {member.user_id === team.owner_id && (
                <span className="text-xs text-neutral-500">(Owner)</span>
              )}
            </li>
          ))}
        </ul>
      </div>

      <p className="w-full text-wrap text-xs text-neutral-500">
        * Invite your team to join by having them submit a request on the
        Explore Teams page.
      </p>
    </div>
  );
}
