import { useFormErrors } from "@/components/Form";
import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import { TextField } from "@/components/ui/TextField";
import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import { useRouter } from "@tanstack/react-router";
import { useState } from "react";
import { DialogTrigger, Form, Heading } from "react-aria-components";
import {
  newTeamSchema,
  useTeamActions,
  type NewTeam,
} from "../hooks/useTeamActions";
import { toast } from "react-toastify";

interface Props {
  eventId: string;
}

export default function NoTeamCard({ eventId }: Props) {
  const router = useRouter();
  const relativeExploreTeamsPath = `/events/${eventId}/dashboard/teams-explorer`;
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
  const { create } = useTeamActions(eventId);

  const form = useForm<
    NewTeam,
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    FormValidateOrFn<NewTeam>,
    undefined,
    undefined,
    undefined
  >({
    defaultValues: {
      name: "",
    },
    onSubmit: async ({ value }) => {
      await create.mutateAsync(value, {
        onSuccess: () => {
          toast.success(`Created team ${value.name} successfully!`);
        },
        onError: () => {
          toast.error("Someting went wrong creating team. Try again!");
        },
      });
      setIsModalOpen(false);
      form.reset();
    },
    validators: {
      onSubmit: newTeamSchema,
    },
  });

  const errors = useFormErrors(form);

  return (
    <div className="border border-input-border rounded-md p-6 w-full md:w-fit md:min-w-120 text-center flex flex-col items-center gap-4">
      <Heading className="text-xl font-semibold text-text-main">
        You&apos;re not in a team yet
      </Heading>
      <p className="text-text-secondary text-base max-w-86">
        Join an existing team to collaborate with others or start your own team
        and invite members.
      </p>

      <div className="flex flex-col sm:flex-row gap-3 mt-2">
          <Button onClick={() => router.navigate({ to: relativeExploreTeamsPath })} variant="secondary">Explore Teams</Button>

        <DialogTrigger isOpen={isModalOpen} onOpenChange={setIsModalOpen}>
          <Button variant="primary">Create Team</Button>

          <Modal>
            <Form
              onSubmit={(e) => {
                e.preventDefault();
                e.stopPropagation();
                form.handleSubmit();
              }}
              validationErrors={errors}
            >
              <div className="flex flex-col gap-4">
                <form.Field name="name">
                  {(field) => (
                    <TextField
                      label="Team Name"
                      name={field.name}
                      placeholder="The Gators"
                      isDisabled={
                        form.state.isSubmitting &&
                        field.state.meta.errors.length === 0
                      }
                      value={field.state.value}
                      onChange={(value) => field.handleChange(value)}
                      validationBehavior="aria"
                    />
                  )}
                </form.Field>

                <Button
                  isPending={form.state.isSubmitting}
                  type="submit"
                  variant="primary"
                >
                  Create
                </Button>
              </div>
            </Form>
          </Modal>
        </DialogTrigger>
      </div>
    </div>
  );
}
