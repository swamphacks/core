import { Button } from "@/components/ui/Button";
import TablerArrowRight from "~icons/tabler/arrow-right";
import { useTeamByInviteId } from "./hooks/useTeamByInviteId";
import { useTeamMembers } from "./hooks/useTeamMembers";
import { useTeamActions } from "./hooks/useTeamActions";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import { useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

interface JoinTeamProps {
  inviteId: string;
}

export default function JoinTeam({ inviteId }: JoinTeamProps) {
  const { data: team, error, isPending, isError } = useTeamByInviteId(inviteId);
  const members = useTeamMembers(team?.id);
  const { joinTeam } = useTeamActions();
  const navigate = useNavigate();

  useEffect(() => {
    async function checkError() {
      if (error instanceof HTTPError) {
        toast.error((await error.response.json()).detail);
      }
    }

    checkError();
  }, [error]);

  if (isError) {
    return (
      <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
        <p>Error loading team</p>
      </div>
    );
  }

  if (isPending || members.isPending) {
    return (
      <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
        <p>Loading team information...</p>
      </div>
    );
  }

  if (!team || !members.data) {
    return (
      <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
        <p>Unable to load team.</p>
      </div>
    );
  }

  const handleJoin = () => {
    joinTeam.mutate(inviteId, {
      onSuccess: () => {
        navigate({ to: "/application" });
      },
      onError: async (error) => {
        if (error instanceof HTTPError) {
          toast.error((await error.response.json()).detail);
        } else {
          toast.error("An error has occurred");
        }
      },
    });
  };

  return (
    <div className="w-full h-full sm:max-w-180 mx-auto font-figtree p-2 relative flex justify-center items-center pb-50">
      <div className="min-w-[35%]">
        <p className="text-text-secondary">Join team</p>
        <p className="text-2xl">{team.name}</p>

        <div className="mt-3">
          <ul className="mt-1 list-disc space-y-1 text-sm text-text-main">
            {members.data.map((member) => (
              <li
                key={member.id}
                className="flex gap-2 items-center mt-2 justify-between"
              >
                <div className="flex gap-2 items-center">
                  <img
                    src={member.image!}
                    alt={"user avatar"}
                    className="size-8 rounded-full object-cover"
                  />
                  <p className="truncate w-40">{member.name}</p>
                </div>

                <span className="text-text-secondary">
                  {team.ownerId === member.id ? "Owner" : "Member"}
                </span>
              </li>
            ))}
          </ul>
        </div>

        <div className="mt-5">
          <Button onClick={handleJoin} className="w-full gap-1">
            Join <TablerArrowRight />
          </Button>
          <p className="text-text-secondary text-xs mt-1">
            To cancel, close this tab.
          </p>
        </div>
      </div>
    </div>
  );
}
