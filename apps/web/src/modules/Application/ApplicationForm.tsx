import { build } from "@/modules/FormBuilder/build";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { QuestionTypes } from "@/modules/FormBuilder/types";
import { showToast } from "@/lib/toast/toast";
import TablerCircleCheck from "~icons/tabler/circle-check";
import { api } from "@/lib/ky";
import { Spinner } from "@/components/ui/Spinner";
import { useMyApplication } from "@/modules/Application/hooks/useMyApplication";
import { formatDistanceToNowStrict, parseISO } from "date-fns";

import Cloud from "./assets/cloud.svg?react";
import Cloud2 from "./assets/cloud2.svg?react";
import Cloud3 from "./assets/cloud3.svg?react";
import Cloud4 from "./assets/cloud4.svg?react";
import Tower from "./assets/tower.svg?react";
import Bell from "./assets/bell.svg?react";

// TODO: dynamically fetch application json data from somewhere (backend, cdn?) instead of hardcoding it in the frontend
import data from "./application.json";

const SAVE_DELAY_MS = 3000; // delay in time before saving form progress

export function ApplicationForm() {
  // TODO: make the `build` api better so components that use this function doesn't have to call useMemo on it?
  const { Form, fieldsTypes } = useMemo(() => build(data), []);
  const fileFields = useRef(new Set<string>());
  const [isInvalid, setIsInvalid] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSavedAt, setLastSavedAt] = useState<Date | undefined>(undefined);
  const [savedText, setSavedText] = useState<string | undefined>("");

  const application = useMyApplication();

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
    if (!application || application.isLoading) return;

    if (application.data?.savedAt) {
      const parsed = parseISO(application.data.savedAt);
      setLastSavedAt(parsed);
    } else {
      setLastSavedAt(undefined);
    }
  }, [application?.data?.savedAt, application?.isLoading]);

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

    const res = await api.post(`application/submit`, {
      body: formData,
    });

    if (!res.ok) {
      const resBody: any = await res.json();

      showToast({
        title: "Submission Error",
        message: resBody.message || "Something went wrong",
        type: "error",
      });

      setIsInvalid(true);
    } else {
      setIsSubmitted(true);
      setIsInvalid(false);
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

  if (application.isLoading) {
    return (
      <div className="flex w-full justify-center pt-10 gap-2 text-text-secondary">
        <Spinner />
        <p>Loading form...</p>
      </div>
    );
  }

  if (!application.data) {
    throw new Error("Application data is empty.");
  }

  const isApplicationSubmitted = application.data.status !== "started";

  const saveStatus = (
    <>
      {isSaving && !isSubmitted && <span>Autosaving...</span>}
      {!isSaving && savedText && !isSubmitted && !isApplicationSubmitted && (
        <span>{savedText}</span>
      )}
    </>
  );

  return (
    <>
      <div className="flex-col">
        <Form
          defaultValues={
            isApplicationSubmitted
              ? undefined
              : JSON.parse(atob(application.data.application))
          }
          onSubmit={onSubmit}
          onNewAttachments={onNewAttachments}
          onChangeDelayMs={SAVE_DELAY_MS}
          onChange={onChange}
          SubmitSuccessComponent={() => (
            <SubmitSuccess submittedAt={application.data.savedAt} />
          )}
          isInvalid={isInvalid}
          isSubmitted={isSubmitted || isApplicationSubmitted}
          isSubmitting={isSubmitting}
          renderFormHeader={(metadata) => {
            return (
              <div className="space-y-3 py-3 rounded-md relative overflow-hidden bg-[#ebf7fc] dark:bg-gray-700 px-4 mt-1">
                <div className="opacity-85">
                  <div className="absolute -bottom-50 -right-33 z-10">
                    <div className="relative inline-block">
                      <Tower className="relative size-90 [transform:rotateX(25deg)_scale(1,0.9)] [transform-origin:bottom_center] z-20" />
                      <Bell className="absolute top-20 right-39 z-10 size-8" />
                      <Cloud className="absolute -top-4 right-20 z-10 size-20 opacity-70 sm:opacity-100" />
                      <Cloud2 className="absolute top-10 right-47 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud3 className="absolute top-1 right-32 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud4 className="absolute -top-5 right-55 z-10 size-20 opacity-30 sm:opacity-100" />
                    </div>
                  </div>
                </div>

                <p className="relative text-2xl text-text-main font-medium z-11 -top-1">
                  {metadata.title}
                </p>
                <p className="relative text-text-main z-11 w-[85%] -top-1 font-medium sm:font-normal">
                  {metadata.description}
                </p>
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
      <div className="flex items-center gap-2 font-medium bg-badge-bg-accepted/50 text-badge-text-accepted rounded-md py-3 pl-3">
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
    </div>
  );
}
