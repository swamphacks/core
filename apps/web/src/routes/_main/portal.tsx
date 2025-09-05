import { EventCard } from "@/features/Event/components/EventCard";
import { useEventsWithUserInfo } from "@/features/Event/hooks/useEventsWithUserInfo";
import { createFileRoute } from "@tanstack/react-router";
import {
  Heading,
  Text,
  ModalOverlay,
  Modal as RAC_Modal,
  Dialog,
  Button,
} from "react-aria-components";
import { useEffect, useState } from "react";
import { TextField } from "@/components/ui/TextField";
import { api } from "@/lib/ky";
import { showToast } from "@/lib/toast/toast";

export const Route = createFileRoute("/_main/portal")({
  component: RouteComponent,
});

function RouteComponent() {
  const { user } = Route.useRouteContext();
  const { data, isLoading, isError } = useEventsWithUserInfo();
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [emailError, setEmailError] = useState("");
  const [preferredName, setPreferredName] = useState("");
  const [preferredEmail, setPreferredEmail] = useState("");

  useEffect(() => {
    const hasSkipped = document.cookie.includes("welcome-modal-skipped=true");
    if (!hasSkipped && user != null && !user.onboarded) {
      setIsModalOpen(true);
    }
  }, []);

  // Cookie helper functions
  const setCookie = (name: string, value: string, days: number = 365) => {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
  };

  const handleSkip = () => {
    setCookie("welcome-modal-skipped", "true");
    setIsModalOpen(false);
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
      // Update user with new information, and set onboarded to true
      await api.patch("users/me", {
        json: {
          name: preferredName || undefined,
          email: preferredEmail || undefined,
        },
      });

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

      console.log("Profile updated successfully:", {
        preferredName,
        preferredEmail,
      });
      setIsModalOpen(false);
    } catch (error) {
      console.error("Failed to update profile:", error);

      showToast({
        title: "Update Failed",
        message: "Failed to update your profile. Please try again.",
        type: "error",
      });
    }
  };

  const modalContent = (
    <div className="p-6 flex flex-col gap-4 w-full max-w-md">
      <div className="text-center">
        <h2 className="text-xl font-semibold mb-2">Welcome!!</h2>
        <p>Please Input Your Preferred Name and Email.</p>
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
          placeholder="Ex: albertgator@gmail.com"
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
          className="px-4 py-2 bg-violet-600 text-white rounded hover:bg-violet-700 transition-colors disabled:text-gray-600 disabled:hover:text-gray-800 disabled:bg-transparent disabled:hover:bg-transparent"
        >
          Submit
        </Button>
      </div>
    </div>
  );

  if (isLoading) {
    return (
      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          {[...Array(3)].map((_, i) => (
            <div
              key={i}
              className="w-96 h-64 bg-neutral-200 dark:bg-neutral-800 animate-pulse rounded"
            />
          ))}
        </div>
      </div>
    );
  }

  if (isError || data === undefined) {
    return (
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          <Text className="text-red-500">
            Whoops, something went wrong, please refresh and try again!
          </Text>
        </div>
      </div>
    );
  }

  return (
    <div>
      <div className="flex flex-col gap-8">
        <div className="flex flex-col gap-2">
          <Heading className="text-3xl text-text-main">
            Welcome, {user?.name ?? "hacker"}!
          </Heading>
          <h2 className="text-text-secondary text-xl">
            Ready to start hacking?
          </h2>
        </div>

        <div className="flex flex-row gap-6">
          {data.map((event) => (
            <EventCard key={event.eventId} {...event} />
          ))}

          {data.length === 0 && (
            <Text className="text-violet-600">
              Awww such empty. Please check back later for events!
            </Text>
          )}
        </div>
      </div>

      <ModalOverlay
        isOpen={isModalOpen}
        onOpenChange={setIsModalOpen}
        className="fixed inset-0 z-50 bg-black/15 backdrop-blur-sm overflow-y-auto flex min-h-full items-center justify-center p-4 text-center"
      >
        <RAC_Modal className="w-full max-w-md max-h-screen sm:max-h-[90vh] overflow-hidden flex flex-col bg-surface text-left align-middle shadow-xs rounded-md">
          <Dialog className="outline-none relative overflow-y-auto">
            {modalContent}
          </Dialog>
        </RAC_Modal>
      </ModalOverlay>
    </div>
  );
}
