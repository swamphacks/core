import { Heading, Text, Form } from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import { TextField } from "@/components/ui/TextField";
import { Switch } from "@/components/ui/Switch";
import IcBaselineDiscord from "~icons/ic/baseline-discord";
import { Button } from "@/components/ui/Button";
import { cn } from "@/utils/cn";
import z from "zod";
import { useFormErrors } from "@/components/Form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import { queryKey as authQueryKey } from "@/lib/auth/hooks/useUser";
import { auth } from "@/lib/authClient";
import { showToast } from "@/lib/toast/toast";
import { useState } from "react";
import TablerPlaystationX from "~icons/tabler/playstation-x";
import TablerLock from "~icons/tabler/lock";

const schema = z.object({
  name: z.string().min(1, "Name is required"),
  preferredEmail: z.preprocess(
    (val: string) => (val === "" ? undefined : val),
    z.email("Email address is invalid").optional(),
  ),
});

export function SettingsPage({
  isOpen,
  toggleSettings,
}: {
  isOpen: boolean;
  toggleSettings: () => void;
}) {
  const queryClient = useQueryClient();
  const { data } = auth.useUser();
  const { user } = data!;

  const updatePersonalInfoMutation = useMutation({
    mutationFn: (data: z.infer<typeof schema>) => {
      return api
        .patch("users/me", {
          json: {
            name: data.name,
            preferred_email: data.preferredEmail,
          },
        })
        .json();
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: authQueryKey });
    },
  });

  const [emailConsent, setEmailConsent] = useState(user?.emailConsent);

  const form = useForm({
    defaultValues: {
      name: user?.name,
      preferredEmail: user?.preferredEmail,
    },
    validators: {
      onChange: schema,
    },
    onSubmit: async ({ value }) => {
      await updatePersonalInfoMutation.mutateAsync(
        value as z.infer<typeof schema>,
        {
          onSuccess: () => {
            form.reset(value);
          },
        },
      );
    },
  });

  const errors = useFormErrors(form);

  const handleEmailConsentToggle = async (selected: boolean) => {
    try {
      setEmailConsent(selected);
      await api.patch("users/me/email-consent", {
        json: { email_consent: selected },
      });
    } catch {
      showToast({
        title: "Something went wrong :(",
        message: "Unable to save changes.",
        type: "error",
      });
      setEmailConsent(!selected);
    }
  };

  // TODO: figure out how to put this settings into a modal overlay so it only renders whenever the user clicks on settings
  // instead of rendering everytime

  return (
    <div
      className={cn(
        "absolute inset-0 z-40 bg-white dark:bg-background transition-all duration-200 p-6",
        isOpen
          ? "opacity-100 visible pointer-events-auto"
          : "opacity-0 invisible pointer-events-none",
      )}
    >
      <div>
        <div className="sm:w-117 flex flex-col gap-8">
          <div className="flex flex-col gap-2">
            <div className="flex justify-between items-center">
              <Heading className="text-2xl text-text-main">Settings</Heading>
              <TablerPlaystationX
                onClick={toggleSettings}
                className="size-7 hover:cursor-pointer text-text-secondary hover:text-text-main"
              />
            </div>
            <h2 className="text-text-secondary text-lg">
              Manage your account settings and preferences.
            </h2>
          </div>

          <div className="flex flex-col gap-6">
            <div>
              <Text className="text-lg font-medium">Personal</Text>

              <div className="mt-3">
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
                        <div>
                          <TextField
                            label="Name"
                            type="text"
                            name={field.name}
                            value={field.state.value}
                            onBlur={field.handleBlur}
                            onChange={(value) => field.handleChange(value)}
                            placeholder="Ex: Albert Gator"
                            validationBehavior="aria"
                          />
                        </div>
                      )}
                    </form.Field>

                    <form.Field name="preferredEmail">
                      {(field) => (
                        <div>
                          <TextField
                            label="Preferred Email"
                            type="email"
                            name={field.name}
                            value={field.state.value ?? ""}
                            onBlur={field.handleBlur}
                            onChange={(value) => field.handleChange(value)}
                            placeholder="Ex: albert.gator@ufl.edu"
                            validationBehavior="aria"
                            description="Email used for communications. If empty, your Account Email will be used."
                          />
                        </div>
                      )}
                    </form.Field>

                    <div>
                      <TextField
                        label="Account Email"
                        type="email"
                        defaultValue={user?.email}
                        placeholder="Ex: albert.gator@ufl.edu"
                        validationBehavior="aria"
                        isDisabled
                        description="This is the email associated with your Discord
                          account. It is set when you signed in with Discord."
                        icon={() => <TablerLock className="text-text-main" />}
                        iconPlacement="right"
                      />
                    </div>
                  </div>

                  <form.Subscribe
                    selector={(state) => [
                      state.canSubmit,
                      state.isSubmitting,
                      state.isDirty,
                    ]}
                  >
                    {([canSubmit, isSubmitting, isDirty]) =>
                      isDirty ? (
                        <Button
                          type="submit"
                          variant="primary"
                          size="sm"
                          isDisabled={!canSubmit || isSubmitting}
                          className={cn(
                            "relative h-8 mt-3 text-white rounded-sm disabled:bg-neutral-200 dark:disabled:bg-neutral-700 disabled:text-neutral-400",
                          )}
                        >
                          {isSubmitting ? "Saving..." : "Save Changes"}
                        </Button>
                      ) : null
                    }
                  </form.Subscribe>
                </Form>
              </div>
            </div>

            <div>
              <Text className="text-lg font-medium">Email Permissions</Text>
              <Switch
                className="mt-3"
                onChange={handleEmailConsentToggle}
                isSelected={emailConsent}
              >
                I consent to receive promotional communications via email.
              </Switch>
            </div>

            <div className="mb-3">
              <Text className="text-lg font-medium block">Socials</Text>
              <div className="inline-flex items-center gap-2 mt-3 rounded-md border-1 bg-button-secondary p-2 border-input-border">
                <span className="bg-button-primary rounded-md p-1">
                  <IcBaselineDiscord
                    width="1.4em"
                    height="1.4em"
                    className="text-white"
                  />
                </span>
                <span>Connected with Discord</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
