import { AvatarStack } from "@/components/ui/AvatarStack";
import TablerUsers from "~icons/tabler/users";
import type { TeamWithMembers } from "../hooks/useMyTeam";
import { useTeamActions } from "../hooks/useTeamActions";
import { Button } from "@/components/ui/Button";
import { toast } from "react-toastify";

interface Props {
  eventId: string;
  team: TeamWithMembers;
}

export default function MyTeamCard({ eventId, team }: Props) {
  const { leave } = useTeamActions(eventId);

  return (
    <div className="border border-input-border rounded-md px-4 py-3 w-full md:w-fit md:min-w-120 flex flex-col gap-4">
      <h3 className="text-text-main text-2xl">{team.name}</h3>

      <div className="flex flex-row text-text-secondary items-center gap-2">
        <AvatarStack
          avatars={team.members.map((member) => ({
            src: member?.image,
            fallback: member?.name,
          }))}
          max={2}
          size="sm"
        />
        <TablerUsers className="h-6 w-6" />
        <p className="text-lg">{team.members.length} / 4 Members</p>
      </div>

      <Button
        variant="danger"
        onPress={() =>
          leave.mutate(team.id, {
            onSuccess: () => {
              toast.success("Left the team successfully.");
            },
            onError: () => {
              toast.error("Something went wrong, try again.");
            },
          })
        }
      >
        Leave
      </Button>
    </div>
  );
}
