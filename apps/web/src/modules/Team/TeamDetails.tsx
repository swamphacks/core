import type { UserContext } from "@/lib/auth/types";
import { useTeamActions } from "./hooks/useTeamActions";
import type { useMyTeam } from "./hooks/useMyTeam";
import { HTTPError } from "ky";
import { toast } from "react-toastify";
import { Button } from "@/components/ui/Button";
import {
  DialogTrigger,
  Focusable,
  Tooltip,
  TooltipTrigger,
} from "react-aria-components";
import TablerPlus from "~icons/tabler/user-plus";
import { Modal } from "@/components/ui/Modal";
import TablerTrash from "~icons/tabler/trash";
import { useTeamInvitation } from "./hooks/useTeamInvite";
import { useState } from "react";
import TablerLink from "~icons/tabler/link";
import TablerLogout from "~icons/tabler/logout";
import TablerShieldCheck from "~icons/tabler/shield-check";
import TablerShieldQuestion from "~icons/tabler/shield-question";
import TablerShieldX from "~icons/tabler/shield-x";
import { Input } from "@/components/ui/Field";

export default function TeamDetails({
  team,
  user,
  accepted,
}: {
  user: UserContext;
  team: NonNullable<ReturnType<typeof useMyTeam>["data"]>;
  accepted: boolean;
}) {
  const { deleteTeam, leaveTeam, kickMember } = useTeamActions();
  const isOwner = team.ownerId === user.userId;
  const members = [...(team.members ?? [])].sort((a, b) => {
    const priority: Record<string, number> = {
      confirmed: 0,
      accepted: 1,
      not_accepted: 2,
    };

    return (priority[a.status] ?? 3) - (priority[b.status] ?? 3);
  });

  const handleDeleteTeam = () => {
    const confirmed = window.confirm("Are you sure you want to delete team?");

    if (confirmed) {
      deleteTeam.mutate(
        { teamId: team.id },
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

  const handleLeaveTeam = () => {
    const confirmed = window.confirm("Are you sure you want to leave team?");

    if (confirmed) {
      leaveTeam.mutate(team.id, {
        onError: async (error) => {
          if (error instanceof HTTPError) {
            toast.error(
              `${(await error.response.json()).detail}. Try refreshing the page.`,
            );
          } else {
            toast.error("Unknown error occurred");
          }
        },
      });
    }
  };

  const handleKickMember = (memberId: string) => {
    kickMember.mutate(
      { teamId: team.id, memberId },
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
  };

  return (
    <div className="mt-3 space-y-2 rounded-md border border-border/70 p-3">
      <div className="flex justify-between">
        <p className="font-medium text-text-main">Your team</p>
        {isOwner && (
          <div className="space-x-2">
            <DialogTrigger onOpenChange={(isOpen) => console.log(isOpen)}>
              <Button className="h-8" size="sm">
                <TablerPlus />
                Invite
              </Button>

              <Modal size="md" isDismissible>
                <InvitationModal teamId={team.id} />
              </Modal>
            </DialogTrigger>

            <Button
              onClick={handleDeleteTeam}
              className="h-8"
              variant="secondary"
              size="sm"
            >
              <TablerTrash />
              Delete
            </Button>
          </div>
        )}
        {!isOwner && (
          <div>
            <Button
              onClick={handleLeaveTeam}
              className="h-8"
              size="sm"
              variant="secondary"
            >
              <TablerLogout />
              Leave
            </Button>
          </div>
        )}
      </div>
      <p className="text-text-secondary">
        Team name:{" "}
        <span className="font-medium text-text-main">{team.name}</span>
      </p>
      <div>
        <p className="text-sm text-text-secondary">Members:</p>
        {members.length > 0 ? (
          <ul className="mt-1 list-disc space-y-1 text-sm text-text-main">
            {members.map((member) => (
              <li
                key={member.id}
                className="flex gap-2 items-center mt-2 justify-between"
              >
                <TeamMember member={member} accepted={accepted} />

                {isOwner && member.id != user.userId && (
                  <Button
                    onClick={() => handleKickMember(member.id)}
                    className="h-8"
                    size="sm"
                    variant="secondary"
                  >
                    Kick
                  </Button>
                )}
              </li>
            ))}
          </ul>
        ) : (
          <p className="mt-1 text-sm text-text-secondary">No members yet.</p>
        )}
      </div>
    </div>
  );
}

function TeamMember({
  member,
  accepted,
}: {
  member: NonNullable<
    NonNullable<ReturnType<typeof useMyTeam>["data"]>["members"]
  >[number];
  accepted: boolean;
}) {
  const getIconForStatus = (status: string) => {
    let icon, tooltip;
    if (status === "confirmed") {
      icon = <TablerShieldCheck className="text-green-500" />;
      tooltip = "Confirmed their attendance 😎";
    } else if (status === "accepted") {
      icon = <TablerShieldQuestion className="text-blue-500" />;
      tooltip = "Accepted, but has not yet confirmed 🙃";
    } else {
      icon = <TablerShieldX />;
      tooltip = "Not yet accepted 🥲";
    }

    return (
      <TooltipTrigger delay={250} closeDelay={250}>
        <Focusable>{icon}</Focusable>
        <Tooltip
          offset={5}
          className="bg-surface border-input-border border-2 py-1 px-2 rounded-md"
        >
          {tooltip}
        </Tooltip>
      </TooltipTrigger>
    );
  };

  return (
    <div className="flex gap-2 items-center">
      <img
        src={member.image!}
        alt={"user avatar"}
        className="size-8 rounded-full object-cover"
      />
      <div className="flex w-50 items-center gap-1">
        <p className="min-w-0 truncate">{member.name}</p>

        {accepted && getIconForStatus(member.status)}
      </div>
    </div>
  );
}

function InvitationModal({ teamId }: { teamId: string }) {
  const { data, isPending } = useTeamInvitation(teamId);
  const [isCopied, setIsCopied] = useState(false);

  const INVITE_BASE_URL = `${window.location.origin}/team/join/`;

  const handleCopy = async () => {
    if (isPending) return;

    try {
      await navigator.clipboard.writeText(INVITE_BASE_URL + data);
      setIsCopied(true);

      setTimeout(() => setIsCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy text: ", err);
    }
  };

  if (isPending) {
    return <p>Loading invitation...</p>;
  }

  return (
    <div className="space-y-3">
      <p>
        Send this private invite link to your teammates. Anyone with this link
        can join your team, so keep it private.
      </p>

      <Input className="rounded-md" value={INVITE_BASE_URL + data} disabled />

      <div className="flex justify-between items-center">
        <Button onClick={handleCopy} size="sm">
          <TablerLink />
          {isCopied ? "Copied!" : "Copy link"}
        </Button>
        <Button slot="close" variant="secondary" size="sm">
          Close
        </Button>
      </div>
    </div>
  );
}
