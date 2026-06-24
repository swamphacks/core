import { build } from "@/modules/FormBuilder/build";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { QuestionTypes } from "@/modules/FormBuilder/types";
import { showToast } from "@/lib/toast/toast";
import TablerCircleCheck from "~icons/tabler/circle-check";
import TablerUpload from "~icons/tabler/upload";
import { api } from "@/lib/ky";
import { Spinner } from "@/components/ui/Spinner";
import { Button } from "@/components/ui/Button";
import { useReplaceResume } from "@/modules/Application/hooks/useReplaceResume";
import { formatDistanceToNowStrict, parseISO } from "date-fns";

// import Cloud from "./assets/cloud.svg?react";
// import Cloud2 from "./assets/cloud2.svg?react";
// import Cloud3 from "./assets/cloud3.svg?react";
// import Cloud4 from "./assets/cloud4.svg?react";
// import Tower from "./assets/tower.svg?react";
// import Bell from "./assets/bell.svg?react";
import Sticker from "./assets/XII_Sticker.svg?react";

// TODO: dynamically fetch application json data from somewhere (backend, cdn?) instead of hardcoding it in the frontend
import data from "./application.json";
import { HTTPError } from "ky";
import type { Hackathon } from "@/modules/Hackathon/hooks/useHackathon";
import type { Application } from "@/modules/Application/hooks/useApplication";

const SAVE_DELAY_MS = 3000; // delay in time before saving form progress

interface ApplicationFormProps {
  hackathon: Hackathon;
  application: Application;
  applicationResponses: any;
}

export function ApplicationForm({
  hackathon,
  application,
  applicationResponses,
}: ApplicationFormProps) {
  // TODO: make the `build` api better so components that use this function doesn't have to call useMemo on it?
  const { Form, fieldsTypes } = useMemo(() => build(data), []);
  const fileFields = useRef(new Set<string>());
  const [isInvalid, setIsInvalid] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSavedAt, setLastSavedAt] = useState<Date | undefined>(undefined);
  const [savedText, setSavedText] = useState<string | undefined>("");
  const [submittedAt, setSubmittedAt] = useState<string | undefined>(undefined);

  // Update saved text every second. Restart interval when lastSavedAt changes.
  useEffect(() => {
    const id = setInterval(() => {
      if (lastSavedAt) {
        setSavedText(
          `Saved ${formatDistanceToNowStrict(lastSavedAt, { addSuffix: true })}`,
        );
      }
    }, 1000);

    return () => clearInterval(id);
  }, [lastSavedAt]);

  // Update saved at status message, if any.
  useEffect(() => {
    if (!application) return;

    if (application?.savedAt) {
      const parsed = parseISO(application.savedAt);
      setLastSavedAt(parsed);
    } else {
      setLastSavedAt(undefined);
    }
  }, [application?.savedAt]);

  const onSubmit = useCallback(async (data: Record<string, any>) => {
    setIsSubmitting(true);

    const formData = new FormData();

    for (const key in data) {
      if (fieldsTypes[key] === QuestionTypes.upload) {
        for (const file of data[key]) {
          formData.append(`${key}[]`, file);
        }
      } else {
        formData.set(key, data[key]);
      }
    }

    try {
      const res = await api.post(`application/submit`, {
        body: formData,
      });

      const submissionResult = await res.json<{ submittedAt: string }>();

      setIsSubmitted(true);
      setIsInvalid(false);
      setSubmittedAt(submissionResult.submittedAt);
    } catch (err) {
      let message = "Something went wrong";
      if (err instanceof HTTPError) {
        const resBody = await err.response.json();
        message = resBody.detail || message;
      }
      showToast({
        title: "Submission Error",
        message: message,
        type: "error",
      });
      setIsInvalid(true);
    }

    setIsSubmitting(false);
  }, []);

  const onNewAttachments = useCallback((newFiles: Record<string, File[]>) => {
    for (const field in newFiles) {
      fileFields.current.add(field);
    }
  }, []);

  const onChange = useCallback(
    async (formValues: Record<string, any>) => {
      if (isSubmitted || isSubmitting) return;

      setIsSaving(true);

      // don't save the file uploads
      for (const field of fileFields.current) {
        delete formValues[field];
      }

      await api.post(`application/save`, {
        json: formValues,
      });

      const now = new Date();

      setLastSavedAt(now);
      // Set text immediately to avoid text flash
      setSavedText(
        `Saved ${formatDistanceToNowStrict(now, { addSuffix: true })}`,
      );
      setIsSaving(false);
    },
    [isSubmitted, isSubmitting],
  );

  const isApplicationSubmitted = application.status !== "started";

  const saveStatus = (
    <>
      {isSaving && !isSubmitted && <span>Autosaving...</span>}
      {!isSaving && savedText && !isSubmitted && !isApplicationSubmitted && (
        <span>{savedText}</span>
      )}
    </>
  );

  const now = new Date();
  const isAccessingEarlyApplication =
    hackathon.acceptEarlyApplications &&
    now >= new Date(hackathon.earlyApplicationOpen as string) &&
    now <= new Date(hackathon.earlyApplicationClose as string);

  return (
    <>
      <div className="flex-col">
        <Form
          defaultValues={
            isApplicationSubmitted ? undefined : applicationResponses
          }
          onSubmit={onSubmit}
          onNewAttachments={onNewAttachments}
          onChangeDelayMs={SAVE_DELAY_MS}
          onChange={onChange}
          SubmitSuccessComponent={() => (
            <SubmitSuccess
              submittedAt={submittedAt || application.submittedAt!}
            />
          )}
          isInvalid={isInvalid}
          isSubmitted={isSubmitted || isApplicationSubmitted}
          isSubmitting={isSubmitting}
          renderFormHeader={(metadata) => {
            return (
              // bg-[#ebf7fc] dark:bg-gray-700
              <div className="space-y-3 py-3 rounded-md relative overflow-hidden mt-1">
                <div className="invisible opacity-65 sm:visible">
                  <div className="absolute right-0 z-10">
                    <Sticker className="size-30 z-1" />
                  </div>
                  {/* <div className="absolute -bottom-50 -right-33 z-10">
                    <div className="relative inline-block">
                      <Tower className="relative size-90 [transform:rotateX(25deg)_scale(1,0.9)] [transform-origin:bottom_center] z-20" />
                      <Bell className="absolute top-20 right-39 z-10 size-8" />
                      <Cloud className="absolute -top-4 right-20 z-10 size-20 opacity-70 sm:opacity-100" />
                      <Cloud2 className="absolute top-10 right-47 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud3 className="absolute top-1 right-32 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud4 className="absolute -top-5 right-55 z-10 size-20 opacity-30 sm:opacity-100" />
                    </div>
                  </div> */}
                </div>

                <p className="relative text-2xl text-text-main font-bold z-11 -top-1">
                  {isAccessingEarlyApplication
                    ? metadata.earlyTitle
                    : metadata.title}
                </p>
                <div className="space-y-2 z-20">
                  <p className="relative z-20">
                    Applications are due{" "}
                    <b>September 1, 2026 at 11:59 PM EST.</b>
                  </p>
                  <p>
                    Questions? Email{" "}
                    <a
                      className="text-text-link"
                      href="mailto:contact@swamphacks.com"
                    >
                      contact@swamphacks.com
                    </a>
                  </p>
                </div>
                {/* <p className="relative text-text-main z-11 w-[85%] -top-1 font-medium sm:font-normal">
                  {metadata.description}
                </p> */}
              </div>
            );
          }}
        />
        <div className="w-full sm:max-w-180 mx-auto font-figtree p-2 text-sm pb-20">
          <div className="lg:hidden">{saveStatus}</div>
        </div>
      </div>
      <div className="hidden fixed xl:inline-flex w-fit z-[999] bottom-3 left-67 text-xs">
        {saveStatus}
      </div>
    </>
  );
}

function SubmitSuccess({ submittedAt }: { submittedAt: string }) {
  return (
    <div>
      <div className="flex items-center gap-2">
        <TablerCircleCheck />
        <p>Thank you! Your application has been received.</p>
      </div>
      <p className="text-gray-500 mt-2">
        Submitted at:{" "}
        {new Intl.DateTimeFormat("en-US", {
          year: "numeric",
          month: "numeric",
          day: "numeric",
          hour: "numeric",
          minute: "2-digit",
        }).format(new Date(submittedAt))}
      </p>
      <ReplaceResume />
    </div>
  );
}

function ReplaceResume() {
  const inputRef = useRef<HTMLInputElement>(null);
  const { mutate: replaceResume, isPending } = useReplaceResume();

  const onSelectFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    // Reset so selecting the same file again still fires onChange
    e.target.value = "";

    if (!file) return;

    if (file.type !== "application/pdf") {
      showToast({
        title: "Invalid file",
        message: "Please upload a PDF resume.",
        type: "error",
      });
      return;
    }

    replaceResume(file, {
      onSuccess: () => {
        showToast({
          title: "Resume updated",
          message: "Your resume has been replaced successfully.",
          type: "success",
        });
      },
      onError: async (err) => {
        let message = "Something went wrong while replacing your resume.";
        if (err instanceof HTTPError) {
          const resBody = await err.response.json<{ detail?: string }>();
          message = resBody.detail || message;
        }

        showToast({
          title: "Upload failed",
          message,
          type: "error",
        });
      },
    });
  };

  return (
    <div className="mt-6 border-t border-gray-200 dark:border-gray-700 pt-4">
      <p className="text-sm text-text-secondary mb-2">
        Uploaded the wrong resume? You can replace it below. This does not
        change any of your other application responses.
      </p>
      <input
        ref={inputRef}
        type="file"
        accept="application/pdf"
        className="hidden"
        onChange={onSelectFile}
      />
      <Button
        variant="secondary"
        size="sm"
        isDisabled={isPending}
        onPress={() => inputRef.current?.click()}
        className="inline-flex items-center gap-2"
      >
        {isPending ? <Spinner /> : <TablerUpload />}
        {isPending ? "Uploading..." : "Replace resume"}
      </Button>
    </div>
  );
}
