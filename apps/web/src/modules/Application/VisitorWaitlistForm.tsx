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

  if (Date.now() >= Date.parse("2026-10-16T04:00:00Z")) {
    return <p>The deadline to join the waitlist has passed.</p>;
  }

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
      <h1 className="text-2xl font-bold">Join the SwampHacks XII Waitlist</h1>
      <p className="my-3 text-text-secondary">
        Submit your information by October 15 at 11:59 PM ET. Joining the
        waitlist does not guarantee admission. If a spot opens, we will email
        you with a confirmation deadline.
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
