import { Heading, Text, Form } from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import { TextField } from "@/components/ui/TextField";
import { Switch } from "@/components/ui/Switch";
import IcBaselineDiscord from "~icons/ic/baseline-discord";
import { Button } from "@/components/ui/Button";
import { cn } from "@/utils/cn";
import z from "zod";
import { useFormErrors } from "@/components/Form";
import { api } from "@/lib/ky";
import { auth } from "@/lib/authClient";
import { showToast } from "@/lib/toast/toast";
import { useState } from "react";
import TablerPlaystationX from "~icons/tabler/playstation-x";
import TablerLock from "~icons/tabler/lock";
import {
  settingsFieldsSchema,
  useSettingsActions,
} from "../hooks/useSettingsActions";
import TablerLogout from "~icons/tabler/logout";
import { useRouter, useCanGoBack } from "@tanstack/react-router";
import TablerHome from "~icons/tabler/home";
import { ThemeSwitch } from "@/components/ThemeProvider";

export function SettingsPage({ logout }: { logout: () => void }) {
  const router = useRouter();
  const canGoBack = useCanGoBack();

  const { data } = auth.useUser();
  const { user } = data!;

  const { updateAccountInfo } = useSettingsActions();

  const [emailConsent, setEmailConsent] = useState(user?.emailConsent);

  const form = useForm({
    defaultValues: {
      name: user?.name,
      preferredEmail: user?.preferredEmail,
    },
    validators: {
      onChange: settingsFieldsSchema,
    },
    onSubmit: async ({ value }) => {
      await updateAccountInfo.mutateAsync(
        value as z.infer<typeof settingsFieldsSchema>,
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

  return (
    <div
      className={cn(
        "absolute inset-0 z-40 bg-white dark:bg-background transition-all duration-200 p-6",
      )}
    >
      <div className="flex justify-center">
        <div className="sm:w-117 flex flex-col gap-8">
          <div className="flex flex-col gap-2">
            <div className="flex justify-between items-center">
              <Heading className="text-2xl text-text-main">Settings</Heading>
              {canGoBack ? (
                <TablerPlaystationX
                  onClick={() => router.history.back()}
                  className="size-7 hover:cursor-pointer text-text-secondary hover:text-text-main"
                />
              ) : (
                <TablerHome
                  onClick={() => router.navigate({ to: "/portal" })}
                  className="size-7 hover:cursor-pointer text-text-secondary hover:text-text-main"
                />
              )}
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
                        description="This email is associated with your third party login. You cannot change this."
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

            <div>
              <Text className="text-lg font-medium">Theme</Text>
              <ThemeSwitch />
            </div>

            <div>
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

            <div className="my-5">
              <Button
                variant="skeleton"
                className="border border-badge-text-rejected text-badge-text-rejected hover:bg-badge-text-rejected/10"
                onClick={logout}
              >
                Log Out <TablerLogout className="text-badge-text-rejected" />
              </Button>
            </div>

            {/* Versioning Footer */}
            <div className="m-4 text-center text-xs text-text-secondary flex flex-col gap-1">
              <a
                href="https://github.com/swamphacks/core"
                target="_blank"
                className="hover:text-text-main"
              >
                beta v1.0.0 - View Source on GitHub
              </a>

              <a
                href="https://discord.gg/hc5RVQnnsU"
                target="_blank"
                className="hover:text-text-main"
              >
                Join us →
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
