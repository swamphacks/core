import { useRouter } from "@tanstack/react-router";
import { build } from "@/features/FormBuilder/build";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { showToast } from "@/lib/toast/toast";
import TablerCircleCheck from "~icons/tabler/circle-check";
import { Link } from "react-aria-components";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import { api } from "@/lib/ky";
import { Spinner } from "@/components/ui/Spinner";
import { useMyApplication } from "@/features/Application/hooks/useMyApplication";
import { formatDistanceToNowStrict, parseISO } from "date-fns";

// TODO: can we put these in the assets folder?
import Cloud from "./cloud.svg?react";
import Cloud2 from "./cloud2.svg?react";
import Cloud3 from "./cloud3.svg?react";
import Cloud4 from "./cloud4.svg?react";
import Tower from "./tower.svg?react";
import Bell from "./bell.svg?react";

// TODO: dynamically fetch application json data from somewhere (backend, cdn?) instead of hardcoding it in the frontend
import data from "@/forms/application.json";

const SAVE_DELAY_MS = 3000; // delay in time before saving form progress

interface ApplicationFormProps {
  eventId: string;
}

export function ApplicationForm({ eventId }: ApplicationFormProps) {
  // TODO: make the `build` api better so components that use this function doesn't have to call useMemo on it?
  const { Form, fieldsMeta } = useMemo(() => build(data), []);
  const fileFields = useRef(new Set<string>());
  const [isInvalid, setIsInvalid] = useState(false);
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const application = useMyApplication(eventId);

  const [isSaving, setIsSaving] = useState(false);
  const [lastSavedAt, setLastSavedAt] = useState<Date | undefined>(undefined);
  const [savedText, setSavedText] = useState<string | undefined>("");

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

  useEffect(() => {
    if (!application || application.isLoading) return;

    if (application.data?.saved_at) {
      const parsed = parseISO(application.data.saved_at);
      setLastSavedAt(parsed);
    } else {
      setLastSavedAt(undefined);
    }
  }, [application?.data?.saved_at, application?.isLoading]);

  const onSubmit = useCallback(async (data: Record<string, any>) => {
    setIsSubmitting(true);

    const formData = new FormData();

    for (const key in data) {
      if (fieldsMeta[key] === QuestionTypes.upload) {
        for (const file of data[key]) {
          formData.append(`${key}[]`, file);
        }
      } else {
        formData.set(key, data[key]);
      }
    }

    const res = await api.post(`events/${eventId}/application/submit`, {
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

      await api.post(`events/${eventId}/application/save`, {
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

  const applicationStatusSubmitted =
    application.data["status"]["application_status"] !== "started";

  const saveStatus = (
    <>
      {isSaving && !isSubmitted && <span>Autosaving...</span>}
      {!isSaving &&
        savedText &&
        !isSubmitted &&
        !applicationStatusSubmitted && <span>{savedText}</span>}
    </>
  );

  return (
    <>
      <div className="hidden fixed xl:inline-flex w-fit z-[999] bottom-3 left-3 text-xs">
        {saveStatus}
      </div>
      <div className="flex-col">
        <Form
          defaultValues={
            applicationStatusSubmitted
              ? undefined
              : JSON.parse(atob(application.data["application"]))
          }
          onSubmit={onSubmit}
          onNewAttachments={onNewAttachments}
          onChangeDelayMs={SAVE_DELAY_MS}
          onChange={onChange}
          SubmitSuccessComponent={SubmitSuccess}
          isInvalid={isInvalid}
          isSubmitted={isSubmitted || applicationStatusSubmitted}
          isSubmitting={isSubmitting}
          // TODO: figure out how to handle this through the JSON file, including how to import the SVG files there as well
          renderFormHeader={(metadata) => {
            return (
              <div className="space-y-3 py-3 rounded-md relative overflow-hidden bg-[#f6fafc] dark:bg-gray-700 px-4 mt-1">
                <div className="opacity-85">
                  <div className="absolute -bottom-50 -right-33 z-10">
                    <div className="relative inline-block">
                      <Tower
                        className="relative size-90 [transform:rotateX(25deg)_scale(1,0.9)]
               [transform-origin:bottom_center] z-20"
                      />
                      <Bell className="absolute top-20 right-39 z-10 size-8" />

                      <Cloud className="absolute -top-4 right-20 z-10 size-20 opacity-70 sm:opacity-100" />
                      <Cloud2 className="absolute top-10 right-47 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud3 className="absolute top-1 right-32 z-10 size-20 opacity-50 sm:opacity-100" />
                      <Cloud4 className="absolute -top-5 right-55 z-10 size-20 opacity-30 sm:opacity-100" />
                    </div>
                  </div>
                </div>

                <p className="relative text-2xl text-text-main font-medium z-50 -top-1">
                  {metadata.title}
                </p>
                <p className="relative text-text-main z-50 w-[85%] -top-1 font-medium sm:font-normal">
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
    </>
  );
}

function SubmitSuccess() {
  const router = useRouter();

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-2 py-3 pl-3 bg-badge-bg-accepted/50 text-badge-text-accepted font-medium">
        <TablerCircleCheck />
        <p>Thank you! Your application has been received.</p>
      </div>
      <Link
        onClick={() =>
          router.navigate({
            to: "/portal",
            replace: true,
            reloadDocument: true,
          })
        }
        className="flex items-center gap-1 text-blue-500 select-none cursor-pointer hover:underline"
      >
        <TablerArrowLeft /> Back to dashboard
      </Link>
    </div>
  );
}
