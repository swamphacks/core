import { Button } from "@/components/ui/Button";
import type { ApplicationFields } from "@/modules/Application/hooks/useApplication";
import { useApplicationForReview } from "@/modules/Application/hooks/useApplicationForReview";
import type { ParsedForm } from "@/modules/Application/hooks/useParsedForm";
import TablerX from "~icons/tabler/x";

interface ApplicationResponsesProps {
  parsedForm: ParsedForm;
  applicationId: string;
}

export default function ApplicationResponsesViewer({
  parsedForm,
  applicationId,
}: ApplicationResponsesProps) {
  const applicationReviewDetails = useApplicationForReview(applicationId);

  if (!applicationReviewDetails.data || applicationReviewDetails.isLoading) {
    return <p>Loading...</p>;
  }

  const appFields = applicationReviewDetails.data.application;
  const resume = applicationReviewDetails.data.resumeUrl;

  const renderAnswer = (value: unknown): string => {
    if (value === null || value === undefined) return "—";
    if (typeof value === "boolean") return value ? "Yes" : "No";
    if (typeof value === "object") return JSON.stringify(value);
    return String(value) === "" ? "-" : String(value);
  };

  const renderFieldResponse = (field: { label: string; name: string }) => {
    let label = field.label;

    if (field.name === "ageCertification") {
      label =
        "I certify that I am 18 years old or will turn 18 before the event date";
    } else if (field.name === "agreeToConduct") {
      label = "I have read and agree to the MLH Code of Conduct";
    } else if (field.name === "infoShareAuthorization") {
      label =
        "I authorize you to share my application/registration information with Major League Hacking. I further agree to the terms of both the MLH Contest Terms and Conditions and the MLH Privacy Policy.";
    }

    return (
      <div
        key={field.name}
        className="rounded-lg border border-input-border bg-card p-2"
      >
        <p className="text-sm font-medium text-text-secondary mb-2">{label}</p>

        <div className="text-base whitespace-pre-wrap break-words leading-relaxed">
          {renderAnswer(appFields[field.name as keyof ApplicationFields])}
        </div>
      </div>
    );
  };

  return (
    <div className="overflow-y-auto max-h-250 p-4">
      <div className="fixed right-5">
        <Button variant="secondary" slot="close">
          <TablerX />
          Close
        </Button>
      </div>

      <div className="mt-2 space-y-6">
        {Object.entries(parsedForm).map(([sectionLabel, fields]) => (
          <div key={sectionLabel}>
            <h3 className="text-lg font-semibold text-text-main mb-4 pb-2 border-b border-input-border">
              {sectionLabel}
            </h3>
            <div className="space-y-3">{fields.map(renderFieldResponse)}</div>
          </div>
        ))}

        <div className="space-y-4 pt-4 border-t border-input-border">
          <h3 className="text-lg font-semibold text-text-main">Resume</h3>
          <div className="p-2 rounded-md border border-input-border h-200 bg-input-bg">
            {resume === "" ? (
              <p>No resume provided.</p>
            ) : (
              <object
                className="w-full h-full"
                type="application/pdf"
                data={resume}
              >
                <p>Unable to load resume.</p>
              </object>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
