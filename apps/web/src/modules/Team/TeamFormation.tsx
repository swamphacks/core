import type { UserContext } from "@/lib/auth/types";
import { HTTPError } from "ky";
import { useState } from "react";
import { toast } from "react-toastify";
import { useTeamActions } from "./hooks/useTeamActions";
import { useMyTeam } from "./hooks/useMyTeam";
import TablerPlus from "~icons/tabler/plus";
import TablerUsersGroup from "~icons/tabler/users-group";
import TablerUserPlus from "~icons/tabler/user-plus";
import { DialogTrigger } from "react-aria-components";
import TeamDetails from "./TeamDetails";
import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { Input } from "@/components/ui/Field";

export default function TeamFormation({ user }: { user: UserContext }) {
  const [teamName, setTeamName] = useState("");
  const { createTeam } = useTeamActions();
  const { data: myTeam, isPending } = useMyTeam();

  const handleCreateTeam = () => {
    if (teamName !== "") {
      createTeam.mutate(
        { name: teamName },
        {
          onError: async (error) => {
            if (error instanceof HTTPError) {
              toast.error((await error.response.json()).detail);
            } else {
              toast.error("Unknown error occurred");
            }
          },
        },
      );
    }
  };

  return (
    <div className="border-t-1 border-border pt-3">
      <div className="flex items-center gap-2">
        <TablerUsersGroup />
        <p>Team Formation</p>
      </div>

      <p className="text-text-secondary mt-2">
        Form a team now to increase your team's chances of being accepted
        together. Although we try our best to keep teams together, we may be
        unable to do so.
      </p>

      {!isPending && myTeam ? (
        <TeamDetails team={myTeam} user={user} />
      ) : (
        <>
          <div className="flex flex-col w-fit">
            <DialogTrigger>
              <Button className="mt-3" size="sm">
                <TablerPlus />
                Create Team
              </Button>

              <Modal size="md" isDismissible>
                <div className="space-y-3">
                  <label>Team Name:</label>
                  <Input
                    value={teamName}
                    onChange={(e) => setTeamName(e.target.value)}
                    className="mt-2 rounded-md"
                    placeholder="Enter a name"
                  />
                  <div className="flex gap-2 justify-end">
                    <Button slot="close" variant="secondary">
                      Cancel
                    </Button>
                    <Button onClick={handleCreateTeam}>Done</Button>
                  </div>
                </div>
              </Modal>
            </DialogTrigger>

            <DialogTrigger>
              <Button variant="secondary" className="mt-3" size="sm">
                <TablerUserPlus />
                Join Team
              </Button>

              <Modal size="md" isDismissible>
                <div className="space-y-3">
                  <p>
                    Ask the team owner to share an invite link with you. Then,
                    open the link in your browser to join the team.
                  </p>
                  <p>
                    You will be able to view your team on this page when you
                    successfully joined.
                  </p>

                  <Button slot="close" variant="secondary">
                    Close
                  </Button>
                </div>
              </Modal>
            </DialogTrigger>
          </div>

          <div className="mt-2">
            <p className="text-text-secondary">
              Looking for a team? Visit our{" "}
              <a
                target="_blank"
                className="text-text-link"
                href="https://discord.com/invite/NfRPv9JtAG"
              >
                Discord server
              </a>{" "}
              to connect with other hackers.
            </p>
          </div>
        </>
      )}
    </div>
  );
}
