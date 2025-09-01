import { useRouter } from "@tanstack/react-router";
import { build } from "@/features/FormBuilder/build";
import { useCallback, useMemo, useRef, useState } from "react";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { showToast } from "@/lib/toast/toast";
import TablerCircleCheck from "~icons/tabler/circle-check";
import { Link } from "react-aria-components";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import { api } from "@/lib/ky";
import { Spinner } from "@/components/ui/Spinner";
import { useApplication } from "@/features/Application/hooks/useApplication";

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

  const application = useApplication(eventId);
  const [isSubmitted, setIsSubmitted] = useState(false);

  const onSubmit = useCallback(async (data: Record<string, any>) => {
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
  }, []);

  const onNewAttachments = useCallback((newFiles: Record<string, File[]>) => {
    for (const field in newFiles) {
      fileFields.current.add(field);
    }
  }, []);

  const onChange = useCallback(
    (formValues: Record<string, any>) => {
      if (isSubmitted) return;

      // don't save the file uploads
      for (const field of fileFields.current) {
        delete formValues[field];
      }

      api.post(`events/${eventId}/application/save`, {
        json: formValues,
      });
    },
    [isSubmitted],
  );

  if (application.isLoading) {
    return (
      <div className="flex w-full justify-center pt-10 gap-2 text-text-secondary">
        <Spinner />
        <p>Loading form...</p>
      </div>
    );
  }

  return (
    <Form
      defaultValues={
        application.data["submitted"]
          ? undefined
          : JSON.parse(atob(application.data["application"]))
      }
      onSubmit={onSubmit}
      onNewAttachments={onNewAttachments}
      onChangeDelayMs={SAVE_DELAY_MS}
      onChange={onChange}
      SubmitSuccessComponent={SubmitSuccess}
      isInvalid={isInvalid}
      isSubmitted={isSubmitted || application.data["submitted"]}
    />
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
        onClick={() => router.navigate({ to: "/portal", replace: true })}
        className="flex items-center gap-1 text-blue-500 select-none cursor-pointer hover:underline"
      >
        <TablerArrowLeft /> Back to dashboard
      </Link>
    </div>
  );
}
