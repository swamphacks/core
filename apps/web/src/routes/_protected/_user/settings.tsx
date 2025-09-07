import { createFileRoute } from "@tanstack/react-router";
import { Heading, Text, Form } from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import { TextField } from "@/components/ui/TextField";
import { Switch } from "@/components/ui/Switch";
import IcBaselineDiscord from "~icons/ic/baseline-discord";

export const Route = createFileRoute("/_protected/_user/settings")({
  component: RouteComponent,
});

function RouteComponent() {
  const test = Route.useRouteContext();

  console.log(test.user);

  const form = useForm({
    defaultValues: {
      preferredName: "",
      preferredEmail: "",
      accountEmail: "",
    },
    // onSubmit: async ({ value }) => {},
  });

  return (
    <div>
      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <Heading className="text-2xl text-text-main">Settings</Heading>
          <h2 className="text-text-secondary text-lg">
            Manage your account settings and preferences.
          </h2>
        </div>

        <div className="sm:w-117 flex flex-col gap-8">
          <div>
            <Text className="text-lg font-medium">Personal</Text>

            <div className="mt-3">
              <Form
                onSubmit={(e) => {
                  e.preventDefault();
                  e.stopPropagation();
                  form.handleSubmit();
                }}
                className="flex flex-col gap-4"
              >
                <div className="flex flex-col gap-4">
                  <form.Field name="preferredName">
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
                          value={field.state.value}
                          onBlur={field.handleBlur}
                          onChange={(value) => field.handleChange(value)}
                          placeholder="Ex: albert.gator@ufl.edu"
                          validationBehavior="aria"
                          description="Email used for communications. If empty, your Account Email will be used."
                        />
                      </div>
                    )}
                  </form.Field>

                  <form.Field name="accountEmail">
                    {(field) => (
                      <div>
                        <TextField
                          label="Account Email"
                          type="email"
                          name={field.name}
                          value={field.state.value}
                          onBlur={field.handleBlur}
                          onChange={(value) => field.handleChange(value)}
                          placeholder="Ex: albert.gator@ufl.edu"
                          validationBehavior="aria"
                          isDisabled
                          description="This is the email associated with your Discord
                          account. It is set when you signed in with Discord."
                        />
                      </div>
                    )}
                  </form.Field>
                </div>
                {/* 
              <div className="flex gap-3 justify-between">
                <form.Subscribe
                  selector={(state) => [state.canSubmit, state.isSubmitting]}
                >
                  {([canSubmit, isSubmitting]) => (
                    <Button
                      type="submit"
                      variant="primary"
                      isDisabled={!canSubmit || isSubmitting}
                      className="flex-1 text-white rounded-sm disabled:bg-neutral-200 dark:disabled:bg-neutral-700 disabled:text-neutral-400"
                    >
                      {isSubmitting ? "Submitting..." : "Continue"}
                    </Button>
                  )}
                </form.Subscribe>
              </div> */}
              </Form>
            </div>
          </div>

          <div>
            <Text className="text-lg font-medium">Email Permissions</Text>
            <Switch className="mt-3">
              I consent to receive promotional communications via email.
            </Switch>
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
        </div>
      </div>
    </div>
  );
}
