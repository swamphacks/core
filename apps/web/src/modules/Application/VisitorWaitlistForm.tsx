import { useQuery } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { build } from "@/modules/FormBuilder/build";
import { QuestionTypes } from "@/modules/FormBuilder/types";
import { api } from "@/lib/ky";
import { HTTPError } from "ky";
import data from "./visitor-waitlist.json";

export default function VisitorWaitlistForm() {
  const { Form, fieldsTypes } = useMemo(() => build(data), []);
  const [isSubmitting, setSubmitting] = useState(false);
  const [isInvalid, setInvalid] = useState(false);
  const [message, setMessage] = useState("");

  const registrationWindow = useQuery({
    queryKey: ["visitor-registration-window"],
    queryFn: ({ signal }) =>
      api
        .get("application/visitor-registration-window", {
          signal,
        })
        .json<{ open: boolean; dayOf: boolean }>(),
    refetchInterval: 30000,
  });

  if (registrationWindow.isPending) return <p>Loading registration…</p>;
  if (registrationWindow.isError) {
    return (
      <p role="alert">
        Unable to load registration availability. Please try again.
      </p>
    );
  }
  if (!registrationWindow.data.open) {
    return <p>Waitlist registration is currently closed.</p>;
  }
  const dayOf = registrationWindow.data.dayOf;

  const submit = async (values: Record<string, unknown>) => {
    setSubmitting(true);
    setInvalid(false);
    setMessage("");
    const body = new FormData();
    for (const [name, value] of Object.entries(values)) {
      if (fieldsTypes[name] === QuestionTypes.upload) {
        for (const file of (value as File[] | undefined) ?? []) {
          body.append(name + "[]", file);
        }
      } else if (value !== undefined && value !== null) {
        body.set(name, String(value));
      }
    }
    try {
      await api.post("application/visitor-waitlist", { body });
      window.location.assign("/application");
    } catch (error) {
      let detail = "Unable to join the waitlist. Please try again.";
      if (error instanceof HTTPError) {
        const response = await error.response
          .json<{ detail?: string }>()
          .catch(() => null);
        detail = response?.detail ?? detail;
      }
      setInvalid(true);
      setMessage(detail);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="w-full sm:max-w-180 mx-auto p-2">
      <h1 className="text-2xl font-bold">
        {dayOf
          ? "SwampHacks XII Day-of Standby Registration"
          : "Join the SwampHacks XII Waitlist"}
      </h1>
      <p className="my-3 text-text-secondary">
        {dayOf
          ? "Register here, then join the day-of standby line at the venue. Previously registered hackers take priority. Wait for staff to accept you; registration does not guarantee admission."
          : "Submit your information by October 15 at 11:59 PM ET. Joining the waitlist does not guarantee admission. If a spot opens, we will email you with a confirmation deadline."}
      </p>
      {message && (
        <p role="alert" className="text-red-500">
          {message}
        </p>
      )}
      <Form
        onSubmit={submit}
        isSubmitting={isSubmitting}
        isInvalid={isInvalid}
        renderFormHeader={() => null}
      />
    </div>
  );
}
