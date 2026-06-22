import { Button } from "@/components/ui/Button";
import { Select } from "@/components/ui/Select";
import { Tab, TabList, TabPanel, TabPanels, Tabs } from "@/components/ui/Tabs";
import { api } from "@/lib/ky";
import { RatingFields } from "@/modules/Application/ApplicationReview/Workspace";
import { type ApplicationFields } from "@/modules/Application/hooks/useApplication";
import type { ParsedForm } from "@/modules/Application/hooks/useParsedForm";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { TextArea } from "react-aria-components";
import TablerX from "~icons/tabler/x";
import TablerCheck from "~icons/tabler/check";
import { useApplicationReviewActions } from "@/modules/Application/hooks/useApplicationReviewActions";
import {
  extendedApplicationQueryKey,
  useExtendedApplication,
  type ParsedExtendedApplicationResponse,
} from "@/modules/Application/hooks/useExtendedApplication";

interface ApplicationResponsesProps {
  parsedForm: ParsedForm;
  applicationId: string;
}

export default function ApplicationViewer({
  parsedForm,
  applicationId,
}: ApplicationResponsesProps) {
  const extendedApplication = useExtendedApplication(applicationId);

  if (!extendedApplication.data || extendedApplication.isLoading) {
    return <p>Loading...</p>;
  }

  const appFields = extendedApplication.data.application;
  const resume = extendedApplication.data.resumeUrl;

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
    <div className="overflow-y-auto max-h-[100vh] p-4">
      <div className="fixed right-5 z-10">
        <Button variant="secondary" slot="close">
          <TablerX />
          Close
        </Button>
      </div>

      <Tabs>
        <TabList aria-label="Tabs" className="text-base">
          <Tab id="responses">Responses</Tab>
          <Tab id="status-and-reviews">Status and Reviews</Tab>
        </TabList>
        <TabPanels>
          <TabPanel id="responses">
            <div className="mt-2 space-y-6 w-full">
              {Object.entries(parsedForm).map(([sectionLabel, fields]) => (
                <div key={sectionLabel}>
                  <h3 className="text-lg font-semibold text-text-main mb-4 pb-2 border-b border-input-border">
                    {sectionLabel}
                  </h3>
                  <div className="flex flex-col gap-3">
                    {fields.map(renderFieldResponse)}
                  </div>
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
          </TabPanel>
          <TabPanel id="status-and-reviews">
            <div className="space-y-3">
              <ApplicationStatus
                status={extendedApplication.data.status}
                applicationId={extendedApplication.data.id}
                userId={extendedApplication.data.user.id}
              />
              {extendedApplication.data.review && (
                <ApplicationReviews
                  review={extendedApplication.data.review}
                  applicationId={extendedApplication.data.id}
                  userId={extendedApplication.data.user.id}
                />
              )}
              <ApplicationAutoDecision
                request={extendedApplication.data.autoDecisionRequest}
                applicationId={extendedApplication.data.id}
                reviewId={extendedApplication.data.review?.id}
                userId={extendedApplication.data.user.id}
              />
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </div>
  );
}

interface ApplicationReviewsProps {
  applicationId: string;
  userId: string;
  review: NonNullable<ParsedExtendedApplicationResponse["review"]>;
}

function ApplicationReviews({
  review,
  applicationId,
}: ApplicationReviewsProps) {
  const {
    experience,
    passion,
    notes,
    isDirty: isReviewsDirty,
    setExperience,
    setPassion,
    setNotes,
    reset: resetRatings,
  } = useReview(
    review.experienceRating || 0,
    review.passionRating || 0,
    review.notes || "",
  );
  const { updateReview } = useApplicationActionsAdmin(applicationId);
  const handleUpdateReview = async () => {
    if (!isReviewsDirty) return;

    await updateReview.mutateAsync({
      reviewId: review.id,
      experienceRating: experience,
      passionRating: passion,
      notes,
    });
  };

  return (
    <div className="text-lg space-y-5">
      <div>
        <div className="flex items-center gap-3">
          <p>Reviewer:</p>
          <div className="flex items-center gap-2">
            {review.reviewer.image && (
              <img
                src={review.reviewer.image}
                alt={"user avatar"}
                className="h-8 w-8 rounded-full object-cover"
              />
            )}
            <span className="inline-block max-w-40 truncate">
              {review.reviewer.name ? review.reviewer.name : "None"}
            </span>
          </div>
        </div>
        <div className="mt-3">
          <RatingFields
            experience={experience}
            passion={passion}
            onExperienceChange={setExperience}
            onPassionChange={setPassion}
          />
          <div className="w-85 flex flex-row gap-4 justify-between items-center mt-3">
            <p className="text-lg items-start self-start">Notes:</p>
            <div className="flex items-center gap-2 justify-end">
              <div className="bg-input-bg rounded-md p-2 h-fit w-fit mt-3 text-base">
                <TextArea
                  className="focus-none outline-none h-18 w-50"
                  placeholder="Additional notes..."
                  value={notes}
                  onChange={(e) => setNotes(e.target.value)}
                />
              </div>
            </div>
          </div>
        </div>
        <div className="flex gap-2 mt-3">
          <Button
            onClick={resetRatings}
            variant="secondary"
            size="sm"
            isDisabled={!isReviewsDirty}
          >
            Reset Ratings
          </Button>
          <Button
            size="sm"
            isDisabled={!isReviewsDirty}
            onClick={handleUpdateReview}
          >
            Update Reviews
          </Button>
        </div>
      </div>
    </div>
  );
}

interface ApplicationStatusProps {
  status: string;
  applicationId: string;
  userId: string;
}

function ApplicationStatus({
  status: initialStatus,
  applicationId,
}: ApplicationStatusProps) {
  const {
    status,
    isDirty: isStatusDirty,
    setStatus,
  } = useStatus(initialStatus);
  const { updateApplication } = useApplicationActionsAdmin(applicationId);

  const handleUpdateStatus = async () => {
    if (!isStatusDirty) return;

    await updateApplication.mutateAsync({
      status,
    });
  };

  return (
    <div>
      <div className="flex gap-2">
        <span className="mr-2 text-lg">Status:</span>
        <Select
          value={status}
          onChange={(v) => {
            if (v != null) {
              setStatus(v.toString());
            }
          }}
          items={[
            {
              id: "started",
              name: "Started",
            },
            {
              id: "submitted",
              name: "Submitted",
            },
            {
              id: "under_review",
              name: "Under Review",
            },
            {
              id: "accepted",
              name: "Accepted",
            },
            {
              id: "rejected",
              name: "Rejected",
            },
            {
              id: "waitlisted",
              name: "Waitlisted",
            },
            {
              id: "withdrawn",
              name: "Withdrawn",
            },
            {
              id: "confirmed",
              name: "Confirmed",
            },
          ]}
          children={null}
        />
      </div>
      <Button
        onClick={handleUpdateStatus}
        size="sm"
        isDisabled={!isStatusDirty}
        className="mt-3"
      >
        Update Status
      </Button>
    </div>
  );
}

interface ApplicationAutoDecisionProps {
  request: ParsedExtendedApplicationResponse["autoDecisionRequest"];
  reviewId?: string;
  applicationId: string;
  userId: string;
}

function ApplicationAutoDecision({
  request,
  reviewId,
  applicationId,
}: ApplicationAutoDecisionProps) {
  const { requestAutoDecision, deleteAutoDecisionRequest } =
    useApplicationReviewActions(applicationId, reviewId);

  const handleRequestAutoAccept = async () => {
    await requestAutoDecision.mutateAsync({
      applicationId,
      accept: true,
      justification: "",
    });
  };

  const handleRequestAutoReject = async () => {
    await requestAutoDecision.mutateAsync({
      applicationId,
      accept: false,
      justification: "",
    });
  };

  const handleUndoAutoDecision = async () => {
    if (!request) return;

    await deleteAutoDecisionRequest.mutateAsync({
      requestId: request.id,
    });
  };

  return (
    <div className="mt-3">
      {request ? (
        <div className="flex items-center gap-2">
          <div
            className={`flex items-center gap-2 px-3 py-1 rounded-md border text-sm font-medium ${
              request.decision === "auto_accept"
                ? "bg-emerald-50 border-emerald-200 text-emerald-800"
                : "bg-rose-50 border-rose-200 text-rose-800"
            }`}
          >
            {request.decision === "auto_accept" ? <TablerCheck /> : <TablerX />}
            <span>
              {request.decision === "auto_accept"
                ? "Auto Accept Requested"
                : "Auto Reject Requested"}
            </span>
          </div>
          <Button
            variant="secondary"
            className="h-9"
            onClick={handleUndoAutoDecision}
          >
            Undo
          </Button>
        </div>
      ) : (
        <div className="flex gap-2 mt-2 mr-2">
          <Button onClick={handleRequestAutoAccept} size="sm">
            <TablerCheck />
            Auto Accept
          </Button>
          <Button onClick={handleRequestAutoReject} size="sm" variant="danger">
            <TablerX />
            Auto Reject
          </Button>
        </div>
      )}
    </div>
  );
}

async function updateApplicationFn(applicationId: string, status: string) {
  try {
    await api.patch(`application`, {
      json: { applicationId, status },
    });

    return { applicationId, status };
  } catch (err) {
    throw err;
  }
}

async function updateReviewFn({
  reviewId,
  passionRating,
  experienceRating,
  notes,
}: {
  reviewId: string;
  passionRating: number;
  experienceRating: number;
  notes: string;
}) {
  try {
    await api.patch(`application/review`, {
      json: {
        id: reviewId,
        passionRating,
        experienceRating,
        notes,
      },
    });

    return {
      reviewId,
      passionRating,
      experienceRating,
      notes,
    };
  } catch (err) {
    throw err;
  }
}

export function useApplicationActionsAdmin(applicationId: string) {
  const queryClient = useQueryClient();

  const updateReview = useMutation({
    mutationFn: (args: {
      reviewId: string;
      passionRating: number;
      experienceRating: number;
      notes: string;
    }) => updateReviewFn(args),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: extendedApplicationQueryKey(applicationId),
      });
    },
  });

  const updateApplication = useMutation({
    mutationFn: ({ status }: { status: string }) =>
      updateApplicationFn(applicationId, status),
    onSuccess: ({ applicationId }) => {
      queryClient.invalidateQueries({
        queryKey: extendedApplicationQueryKey(applicationId),
      });
    },
  });

  return { updateApplication, updateReview };
}

function useStatus(initialStatus: string) {
  const [status, setStatus] = useState(initialStatus);
  const [isDirty, setDirty] = useState(false);

  useEffect(() => {
    setDirty(status !== initialStatus);
  }, [status, initialStatus]);

  useEffect(() => {
    setStatus(initialStatus);
  }, [initialStatus]);

  const reset = () => {
    setStatus(initialStatus);
  };

  return {
    status,
    setStatus,
    isDirty,
    reset,
  };
}

function useReview(
  initialExperience: number,
  initialPassion: number,
  initialNotes: string,
) {
  const [experience, setExperience] = useState(initialExperience);
  const [passion, setPassion] = useState(initialPassion);
  const [notes, setNotes] = useState(initialNotes);
  const [isDirty, setDirty] = useState(false);

  useEffect(() => {
    setDirty(
      experience !== initialExperience ||
        passion !== initialPassion ||
        notes !== initialNotes,
    );
  }, [
    experience,
    passion,
    notes,
    initialNotes,
    initialExperience,
    initialPassion,
  ]);

  useEffect(() => {
    setExperience(initialExperience);
    setPassion(initialPassion);
    setNotes(initialNotes);
  }, [initialExperience, initialPassion, initialNotes]);

  const reset = () => {
    setExperience(initialExperience);
    setPassion(initialPassion);
    setNotes(initialNotes);
  };

  return {
    experience,
    passion,
    notes,
    isDirty,
    setExperience,
    setPassion,
    setNotes,
    reset,
  };
}
