import {
  ModalOverlay,
  Modal as RAC_Modal,
  Dialog,
  Button,
  Form,
  Heading,
} from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import { z } from "zod";
import { TextField } from "@/components/ui/TextField";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import Cookies from "js-cookie";

interface OnboardingModalProps {
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
}

const onboardingSchema = z.object({
  preferredName: z.string().min(1, "A preferred name is required"),
  preferredEmail: z
    .string()
    .min(1, "A preferred email is required")
    .email("Please enter a valid email address"),
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
        // Update user with new information
        await api.patch("users/me", {
          json: {
            name: value.preferredName,
            email: value.preferredEmail,
          },
        });

        // Mark user as onboarded
        await api.patch("users/me/onboarded", {
          json: {
            onboarded: true,
          },
        });

        showToast({
          title: "Profile Updated",
          message: "Your profile has been updated successfully!",
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
            <div className="text-center">
              <Heading slot="title" className="text-xl font-semibold mb-2">
                Welcome!
              </Heading>
              <p>Please input your preferred name and email.</p>
            </div>
            <div className="flex flex-col gap-2">
              <form.Field name="preferredName">
                {(field) => (
                  <div>
                    <TextField
                      label="Preferred Name"
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

            <div className="flex gap-3 justify-end">
              <Button
                type="button"
                onPress={handleSkip}
                className="px-4 py-2 text-gray-600 hover:text-gray-800 transition-colors"
              >
                Skip
              </Button>
              <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
              >
                {([canSubmit, isSubmitting]) => (
                  <Button
                    type="submit"
                    isDisabled={!canSubmit || isSubmitting}
                    className="px-4 py-2 bg-violet-600 text-white rounded hover:bg-violet-700 transition-colors disabled:opacity-50 disabled:bg-transparent disabled:text-gray-600"
                  >
                    {isSubmitting ? "Submitting..." : "Submit"}
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
