import { useState } from "react";
import {
  ModalOverlay,
  Modal as RAC_Modal,
  Dialog,
  Button,
} from "react-aria-components";
import { TextField } from "@/components/ui/TextField";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";
import Cookies from "js-cookie";

interface OnboardingModalProps {
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
}

export function OnboardingModal({
  isOpen,
  onOpenChange,
}: OnboardingModalProps) {
  const [emailError, setEmailError] = useState("");
  const [preferredName, setPreferredName] = useState("");
  const [preferredEmail, setPreferredEmail] = useState("");

  const closeModal = () => onOpenChange(false);

  const handleSkip = () => {
    Cookies.set("welcome-modal-skipped", "true", { expires: 365 });
    closeModal();
  };

  const handleEmailChange = (value: string) => {
    setPreferredEmail(value);
    const emailRegex = /\S+@\S+\.\S{2,}/;
    if (value && !emailRegex.test(value)) {
      setEmailError(
        "Please enter a valid email address (e.g., user@example.com).",
      );
    } else {
      setEmailError("");
    }
  };

  const handleSubmit = async () => {
    try {
      // Update user with new information
      await api.patch("users/me", {
        json: {
          name: preferredName || undefined,
          email: preferredEmail || undefined,
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
  };

  return (
    <ModalOverlay
      isOpen={isOpen}
      onOpenChange={onOpenChange}
      className="fixed inset-0 z-50 bg-black/15 backdrop-blur-sm overflow-y-auto flex min-h-full items-center justify-center p-4 text-center"
    >
      <RAC_Modal className="w-full max-w-md max-h-screen sm:max-h-[90vh] overflow-hidden flex flex-col bg-surface text-left align-middle shadow-xs rounded-md">
        <Dialog className="outline-none relative overflow-y-auto p-6 flex flex-col gap-4">
          <div className="text-center">
            <h2 className="text-xl font-semibold mb-2">Welcome!</h2>
            <p>Please input your preferred name and email.</p>
          </div>

          <div className="flex flex-col gap-2">
            <TextField
              label="Preferred Name"
              value={preferredName}
              onChange={setPreferredName}
              placeholder="Ex: Albert Gator"
            />
            <TextField
              label="Preferred Email"
              value={preferredEmail}
              onChange={handleEmailChange}
              placeholder="Ex: albert.gator@ufl.edu"
            />
            {emailError && (
              <p className="text-red-200 text-sm mt-1">{emailError}</p>
            )}
          </div>

          <div className="flex gap-3 justify-end">
            <Button
              onPress={handleSkip}
              className="px-4 py-2 text-gray-600 hover:text-gray-800 transition-colors"
            >
              Skip
            </Button>
            <Button
              onPress={handleSubmit}
              isDisabled={!!emailError || !preferredName || !preferredEmail}
              className="px-4 py-2 bg-violet-600 text-white rounded hover:bg-violet-700 transition-colors disabled:opacity-50"
            >
              Submit
            </Button>
          </div>
        </Dialog>
      </RAC_Modal>
    </ModalOverlay>
  );
}
