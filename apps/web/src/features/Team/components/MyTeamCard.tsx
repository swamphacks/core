import { AvatarStack } from "@/components/ui/AvatarStack";
import TablerUsers from "~icons/tabler/users";
import type { TeamWithMembers } from "../hooks/useMyTeam";
import { useTeamActions } from "../hooks/useTeamActions";
import { Button } from "@/components/ui/Button";
import { toast } from "react-toastify";
import TablerDoorExit from "~icons/tabler/door-exit";
import TablerCircleX from "~icons/tabler/circle-x";
import TablerPlus from "~icons/tabler/plus";
import { Tooltip } from "@/components/ui/Tooltip";
import { Modal } from "@/components/ui/Modal";
import { TextField } from "@/components/ui/TextField";
import { DialogTrigger, Form } from "react-aria-components";
import { useState } from "react";
import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import { useFormErrors } from "@/components/Form";
import { api } from "@/lib/ky";
import { HTTPError } from "ky";
import { z } from "zod";

interface Props {
  eventId: string;
  userId: string;
  team: TeamWithMembers;
}

const inviteEmailSchema = z.object({
  email: z.string().email("Please enter a valid email address"),
});

type InviteEmail = z.infer<typeof inviteEmailSchema>;

async function inviteUserToTeam(teamId: string, email: string) {
  try {
    await api.post(`teams/${teamId}/invite`, {
      json: { email },
    });
  } catch (err) {
    if (err instanceof HTTPError) {
      if (err.response.status === 409) {
        toast.error("A pending invitation already exists for this email.");
      } else if (err.response.status === 403) {
        toast.error("You don't have permission to invite users to this team.");
      } else {
        const data = await err.response.json().catch(() => ({}));
        toast.error(data.message || "Failed to send invitation. Try again.");
      }
    }
    throw err;
  }
}

export default function MyTeamCard({ eventId, userId, team }: Props) {
  const { leave, kickTeamMember } = useTeamActions(eventId);
  const [isInviteModalOpen, setIsInviteModalOpen] = useState(false);
  const isOwner = userId === team.owner_id;
  const isTeamFull = team.members.length >= 4;

  const inviteForm = useForm<
    InviteEmail,
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    FormValidateOrFn<InviteEmail>,
    undefined,
    undefined,
    undefined
  >({
    defaultValues: {
      email: "",
    },
    onSubmit: async ({ value }) => {
      await inviteUserToTeam(team.id, value.email);
      toast.success(`Invitation sent to ${value.email}!`);
      setIsInviteModalOpen(false);
      inviteForm.reset();
    },
    validators: {
      onSubmit: inviteEmailSchema,
    },
  });

  const inviteFormErrors = useFormErrors(inviteForm);

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

        <Tooltip
          tooltipProps={{
            label: "Leave Team",
            offset: 4,
          }}
          triggerProps={{
            delay: 150,
          }}
        >
          <Button
            onClick={handleLeaveTeam}
            variant="icon"
            className="aspect-square p-2.5"
          >
            <TablerDoorExit className="w-5 h-5 text-red-600" />
          </Button>
        </Tooltip>
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
        {isOwner && !isTeamFull && (
          <DialogTrigger
            isOpen={isInviteModalOpen}
            onOpenChange={setIsInviteModalOpen}
          >
            <Tooltip
              tooltipProps={{
                label: "Invite Member",
                offset: 4,
              }}
              triggerProps={{
                delay: 150,
              }}
            >
              <Button variant="icon" className="aspect-square p-2.5 ml-auto">
                <TablerPlus className="w-5 h-5 text-green-600" />
              </Button>
            </Tooltip>

            <Modal title="Invite Team Member" size="md">
              <Form
                onSubmit={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  inviteForm.handleSubmit();
                }}
                validationErrors={inviteFormErrors}
              >
                <div className="flex flex-col gap-4">
                  <inviteForm.Field name="email">
                    {(field) => (
                      <TextField
                        label="Email Address"
                        name={field.name}
                        placeholder="teammate@example.com"
                        type="email"
                        isDisabled={
                          inviteForm.state.isSubmitting &&
                          field.state.meta.errors.length === 0
                        }
                        value={field.state.value}
                        onChange={(value) => field.handleChange(value)}
                        validationBehavior="aria"
                        errorMessage={field.state.meta.errors[0]}
                        isRequired
                      />
                    )}
                  </inviteForm.Field>

                  <div className="flex flex-row gap-3 justify-end">
                    <Button
                      variant="secondary"
                      onPress={() => {
                        setIsInviteModalOpen(false);
                        inviteForm.reset();
                      }}
                      isDisabled={inviteForm.state.isSubmitting}
                    >
                      Cancel
                    </Button>
                    <Button
                      isPending={inviteForm.state.isSubmitting}
                      type="submit"
                      variant="primary"
                    >
                      Send Invitation
                    </Button>
                  </div>
                </div>
              </Form>
            </Modal>
          </DialogTrigger>
        )}
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
                <Tooltip
                  tooltipProps={{
                    label: `Remove ${member.name || "member"}`,
                    offset: 4,
                  }}
                  triggerProps={{
                    delay: 150,
                  }}
                >
                  <Button
                    variant="unstyled"
                    className="ml-2 p-1"
                    onClick={() => handleKickMember(member.user_id)}
                  >
                    <TablerCircleX className="w-5 h-5 hover:text-red-600 hover:cursor-pointer transition-colors duration-150" />
                  </Button>
                </Tooltip>
              )}

              {/* Indicate the owner */}
              {member.user_id === team.owner_id && (
                <span className="text-xs text-neutral-500 select-none">
                  (Owner)
                </span>
              )}
            </li>
          ))}
        </ul>
      </div>

      {isOwner && (
        <p className="w-full text-wrap text-xs text-neutral-500">
          {isTeamFull
            ? "* Team is full (4/4 members). Remove a member to invite someone new."
            : "* Click the + button to invite members by email, or have them submit a request on the Explore Teams page."}
        </p>
      )}
      {!isOwner && (
        <p className="w-full text-wrap text-xs text-neutral-500">
          * Ask the team owner to invite members or submit a request on the
          Explore Teams page.
        </p>
      )}
    </div>
  );
}
