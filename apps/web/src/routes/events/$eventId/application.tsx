import {
  createFileRoute,
  useLocation,
  useRouter,
} from "@tanstack/react-router";
import { build } from "@/features/FormBuilder/build";
import data from "@/forms/application.json";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { env } from "@/config/env";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { ErrorBoundary } from "react-error-boundary";
import TablerAlertCircle from "~icons/tabler/alert-circle";
import { showToast } from "@/lib/toast/toast";
import TablerCircleCheck from "~icons/tabler/circle-check";
import { Link } from "react-aria-components";
import TablerArrowLeft from "~icons/tabler/arrow-left";

export const Route = createFileRoute("/events/$eventId/application")({
  component: RouteComponent,
});

function RouteComponent() {
  // const [submitError, setSubmitError] = useState<string | null>(null);

  return (
    <ErrorBoundary FallbackComponent={Fallback}>
      <div className="w-full h-screen bg-surface">
        {/* this padding left prevent the page being shifted when the form fully loads because of the scrollbar */}
        <div className="w-full bg-surface transition-[background] sm:pl-[calc(100vw-100%)]">
          <ApplicationForm />
        </div>
      </div>
    </ErrorBoundary>
  );
}

const SAVE_DELAY_MS = 1000; // delay in time before saving form progress

// TODO: Fetch the approriate form data corresponding to the event
export const ApplicationForm = () => {
  const location = useLocation();
  // TODO: make the `build` api better so components that use this function doesn't have to call useMemo on it?
  const { Form, fieldsMeta } = useMemo(() => build(data), []);
  const [defaultValues, setDefaultValues] = useState<Record<string, any>>({});
  const fileFields = useRef(new Set<string>());
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [isInvalid, setIsInvalid] = useState(false);

  useEffect(() => {
    fetch(`${env.BASE_API_URL}${location.pathname}`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setDefaultValues(JSON.parse(atob(data["application"])));
      });
  }, []);

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

    const res = await fetch(`${env.BASE_API_URL}${location.pathname}/submit`, {
      method: "POST",
      body: formData,
      credentials: "include",
    });

    if (!res.ok) {
      const resBody = await res.json();

      showToast({
        title: "Submission Error",
        message: resBody.message || "Something went wrong",
        type: "error",
      });

      setIsInvalid(true);
    } else {
      setIsSubmitted(true);
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

      fetch(`${env.BASE_API_URL}${location.pathname}/save`, {
        method: "POST",
        body: JSON.stringify(formValues),
        credentials: "include",
      });
    },
    [isSubmitted],
  );

  return (
    <Form
      defaultValues={defaultValues}
      onSubmit={onSubmit}
      onNewAttachments={onNewAttachments}
      onChangeDelayMs={SAVE_DELAY_MS}
      onChange={onChange}
      SubmitSuccessComponent={SubmitSuccess}
      isInvalid={isInvalid}
    />
  );
};

function Fallback() {
  return (
    <div className="w-full h-full bg-surface flex justify-center items-center gap-2 text-red-400">
      <TablerAlertCircle />
      <p>Something went wrong while loading form :(</p>
    </div>
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
