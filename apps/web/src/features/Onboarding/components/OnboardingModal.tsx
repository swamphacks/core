import {
  ModalOverlay,
  Modal as RAC_Modal,
  Dialog,
  Form,
  Heading,
} from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { TextField } from "@/components/ui/TextField";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import Cookies from "js-cookie";
import { Button } from "@/components/ui/Button";

interface OnboardingModalProps {
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
}

const onboardingSchema = z.object({
  preferredName: z.string().min(1, "A preferred name is required"),
  preferredEmail: z.email("Please enter a valid email address"),
});

export function OnboardingModal({
  isOpen,
  onOpenChange,
}: OnboardingModalProps) {
  const form = useForm({
    defaultValues: {
      preferredName: "",
      preferredEmail: "",
    },
    validators: {
      onChange: onboardingSchema,
      onMount: onboardingSchema,
    },
    onSubmit: async ({ value }) => {
      try {
        await api.patch("users/me/onboarding", {
          json: {
            preferred_email: value.preferredEmail,
            name: value.preferredName,
          },
        });

        showToast({
          title: "Profile Updated",
          message: "Your profile has been updated successfully.",
          type: "success",
        });

        closeModal();
      } catch (error) {
        console.error("Failed to update profile:", error);
        showToast({
          title: "Update Failed",
          message: "Failed to update your profile. Please try again.",
          type: "error",
        });
      }
    },
  });

  const closeModal = () => onOpenChange(false);

  const handleSkip = () => {
    Cookies.set("welcome-modal-skipped", "true", { expires: 365 });
    closeModal();
  };

  return (
    <ModalOverlay
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      className="fixed inset-0 z-50 bg-black/15 backdrop-blur-sm overflow-y-auto flex min-h-full items-center justify-center p-4 text-center"
    >
      <RAC_Modal className="w-full max-w-md max-h-screen sm:max-h-[90vh] overflow-hidden flex flex-col bg-surface text-left align-middle shadow-xs rounded-md">
        <Dialog className="outline-none relative overflow-y-auto p-6">
          <Form
            onSubmit={(e) => {
              e.preventDefault();
              e.stopPropagation();
              form.handleSubmit();
            }}
            className="flex flex-col gap-4"
          >
            <div className="text-center flex flex-col justify-center items-center">
              <img
                className="dark:hidden h-24 w-24"
                src="/core-favicon-light.svg"
              />
              <img
                className="dark:block hidden h-24 w-24"
                src="/core-favicon-dark.svg"
              />
              <Heading slot="title" className="text-2xl font-bold mb-2">
                Tell us about yourself!
              </Heading>
            </div>
            <div className="flex flex-col gap-2">
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
                    />
                  </div>
                )}
              </form.Field>
            </div>

            <div className="flex gap-3 justify-between">
              <Button
                variant="skeleton"
                onPress={handleSkip}
                className="flex-1 border-input-border border-[1px] rounded-sm"
              >
                Skip for now
              </Button>
              <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
              >
                {([canSubmit, isSubmitting]) => (
                  <Button
                    type="submit"
                    variant="primary"
                    isDisabled={!canSubmit || isSubmitting}
                    className="bg-cyan-700 flex-1 text-white rounded-sm disabled:bg-neutral-200 dark:disabled:bg-neutral-700 disabled:text-neutral-400 hover:bg-cyan-800 pressed:bg-cyan-700"
                  >
                    {isSubmitting ? "Submitting..." : "Continue"}
                  </Button>
                )}
              </form.Subscribe>
            </div>
          </Form>
        </Dialog>
      </RAC_Modal>
    </ModalOverlay>
  );
}
