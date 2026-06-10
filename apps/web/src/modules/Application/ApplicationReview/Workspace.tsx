import type { UserContext } from "@/lib/auth/types";
import { Link } from "@tanstack/react-router";
import { useRatings } from "@/modules/Application/hooks/useRatings";
import { useAssignedApplications } from "@/modules/Application/hooks/useAssignedApplications";
import { useAppReviewProgress } from "@/modules/Application/hooks/useAppReviewProgress";
import type { AssignedApplications } from "@/lib/openapi/types";
import { useApplicationForReview } from "@/modules/Application/hooks/useApplicationForReview";
import { useApplicationReviewActions } from "@/modules/Application/hooks/useApplicationReviewActions";
import { toast } from "react-toastify";
import { ProgressBar } from "@/components/ui/ProgressBar";
import TablerX from "~icons/tabler/x";
import { useState } from "react";
import { DialogTrigger } from "react-aria-components";
import { Modal } from "@/components/ui/Modal";
import { Input } from "@/components/ui/Field";
import { Rating } from "@smastrom/react-rating";
import TablerCheck from "~icons/tabler/check";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import TablerArrowRight from "~icons/tabler/arrow-right";
import TablerRefresh from "~icons/tabler/refresh";
import { Button } from "@/components/ui/Button";

interface ApplicationReviewWorkspaceProps {
  user: UserContext;
}

export default function ApplicationReviewWorkspace({
  user,
}: ApplicationReviewWorkspaceProps) {
  const assignedApps = useAssignedApplications();
  const appReviewProgress = useAppReviewProgress(assignedApps.data || []);

  if (assignedApps.isLoading) {
    return (
      <div className="flex gap-2">
        <p>Loading assigned applications...</p>
      </div>
    );
  }

  if (!assignedApps.data || assignedApps.data.length === 0) {
    return (
      <div>
        <p className="text-text-secondary">
          You have no assigned applications to review at this time.
          <br />
          Refresh the page or ask the organizers if you believe this is an
          error.
        </p>
      </div>
    );
  }

  if (appReviewProgress.finished) {
    return (
      <div className="p-3 flex justify-center">
        <div className="flex flex-col items-start bg-surface w-full max-w-md px-6 py-6 border border-border rounded-lg shadow-sm">
          {/* Image */}
          <div className="mb-6 w-full flex flex-row justify-center">
            <img
              src="https://dejpknyizje2n.cloudfront.net/media/carstickers/versions/happy-go-lucky-alligator-art-sticker-uf6a4-438d-x450.png"
              alt="Review Completed"
              className="w-48 h-48 object-contain"
            />
          </div>

          <h2 className="text-2xl font-semibold mb-3">Reviews Completed</h2>

          <p className="mb-6 text-text-secondary">
            You have completed reviewing all assigned applications. Thank you
            for your time and effort! You can go back to change your reviews if
            needed.
          </p>

          <div className="flex flex-row gap-4">
            <Button variant="secondary" onClick={appReviewProgress.goPrevious}>
              Change Reviews
            </Button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-2">
      <div className="flex flex-col md:flex-row md:items-center gap-3 text-lg justify-between">
        <div className="flex items-center gap-3">
          <Link to="/">
            <Button variant="secondary" size="sm">
              <TablerArrowLeft />
              Home
            </Button>
          </Link>
          <h1>SwampHacks Application Review</h1>
        </div>

        <div className="pr-3 min-w-100">
          <ProgressBar
            className="w-full"
            label="Review Progress"
            maxValue={appReviewProgress.totalApplications}
            value={appReviewProgress.currentIndex + 1}
            valueLabel={`${appReviewProgress.currentIndex + 1}/${appReviewProgress.totalApplications} applications`}
          />
        </div>
      </div>

      {appReviewProgress.currentAssignedApplication ? (
        <ApplicationViewer
          user={user}
          next={appReviewProgress.goNext}
          back={appReviewProgress.goPrevious}
          currentIndex={appReviewProgress.currentIndex}
          totalApplications={appReviewProgress.totalApplications}
          assignedApplication={appReviewProgress.currentAssignedApplication}
        />
      ) : (
        <div>Loading application...</div>
      )}
    </div>
  );
}

interface ApplicationViewerProps {
  user: UserContext;
  assignedApplication: AssignedApplications[number];
  currentIndex: number;
  totalApplications: number;
  next: () => void;
  back: () => void;
}

function ApplicationViewer({
  user,
  assignedApplication,
  totalApplications,
  currentIndex,
  next,
  back,
}: ApplicationViewerProps) {
  const applicationReviewDetails = useApplicationForReview(
    assignedApplication.applicationId,
  );
  const { review, requestAutoDecision, deleteAutoDecisionRequest } =
    useApplicationReviewActions(assignedApplication.applicationId);
  const { experience, passion, isDirty, setExperience, setPassion, reset } =
    useRatings(
      applicationReviewDetails.data?.experienceRating || 0,
      applicationReviewDetails.data?.passionRating || 0,
    );

  if (!applicationReviewDetails.data || applicationReviewDetails.isLoading) {
    return <p>Loading...</p>;
  }
  const getHackathonExperienceText = (experience: string) => {
    switch (experience) {
      case "first_time":
        return "Swamphacks would be my first!";
      case "one":
        return "1";
      case "two":
        return "2";
      case "three":
        return "3";
      case "four_or_more":
        return "4+";
      default:
        return "";
    }
  };

  const getProjectExperienceText = (experience: string) => {
    switch (experience) {
      case "no_experience":
        return "Swamphacks would be my first!";
      case "course_experience":
        return "From courses";
      case "independent_project":
        return "Yes";
      default:
        return "";
    }
  };

  const handleSubmitReview = async () => {
    if (!allowSubmit) return;

    await review.mutateAsync(
      {
        experienceRating: experience,
        passionRating: passion,
      },
      {
        onSuccess: () => {
          reset();
          next();
        },
        onError: () => {
          toast.error("Failed to submit review. Please try again.");
          reset();
        },
      },
    );
  };

  const handleRequestAutoAccept = async (justification: string) => {
    await requestAutoDecision.mutateAsync({
      applicationId: assignedApplication.applicationId,
      accept: true,
      justification,
    });
  };

  const handleRequestAutoReject = async (justification: string) => {
    await requestAutoDecision.mutateAsync({
      applicationId: assignedApplication.applicationId,
      accept: false,
      justification,
    });
  };

  const handleUndoAutoDecision = async () => {
    if (applicationReviewDetails.data.autoDecisionRequestId === null) return;

    await deleteAutoDecisionRequest.mutateAsync({
      requestId: applicationReviewDetails.data.autoDecisionRequestId,
    });
  };

  const appFields = applicationReviewDetails.data.application;
  const resume = applicationReviewDetails.data.resumeUrl;
  const autoDecision = applicationReviewDetails.data.autoDecision;
  const isCompleted = assignedApplication.status === "completed";
  const isLast = currentIndex === totalApplications - 1;
  const isFilled = experience > 0 && passion > 0;
  const allowSubmit = isFilled && isDirty;

  const mode = (() => {
    if (!isCompleted) return "submit"; // Not last, not submitted yet
    if (isCompleted && isDirty) return "completed-dirty";
    return "completed-clean"; // Not last, already submitted and clean
  })();

  function ApplicantInfo() {
    return (
      <div className="p-4 rounded-md border border-input-border bg-card">
        <h2 className="block mb-3 text-lg font-medium">
          Applicant Information
        </h2>
        <div className="flex gap-15 items-start">
          <div className="grid grid-cols-1 gap-4">
            <div>
              <div className="text-text-secondary">Name</div>
              <div className="font-medium">
                {appFields.firstName + " " + appFields.lastName}
              </div>
            </div>

            <div>
              <div className="text-text-secondary">Major(s)</div>
              <div className="font-medium">{appFields.majors}</div>
            </div>

            <div>
              <div className="text-text-secondary">School</div>
              <div className="truncate max-w-60">{appFields.school}</div>
            </div>

            <div>
              <div className="text-text-secondary">Graduation Year</div>
              <div className="font-medium">{appFields.graduationYear}</div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4">
            <div>
              <div className="text-text-secondary">
                # of Hackathons Attended
              </div>
              <div className="font-medium">
                {getHackathonExperienceText(appFields.experience)}
              </div>
            </div>
            <div>
              <div className="text-text-secondary">Project Experience</div>
              <div className="font-medium">
                {getProjectExperienceText(appFields.projectExperience)}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  function Essays() {
    return (
      <div className="space-y-4">
        <div className="p-4 rounded-md border border-input-border bg-card">
          <h2 className="block mb-3 text-lg font-medium">Essay Responses</h2>
          <div className="space-y-6">
            <div>
              <div className="text-text-secondary mb-2">
                What is your most memorable experience working in a group? What
                did you learn and accomplish?
              </div>
              <div className="whitespace-pre-wrap bg-surface p-3 rounded-md">
                {appFields.essay1}
              </div>
            </div>

            <div>
              <div className="text-text-secondary mb-2">
                Tell us about a project you are most proud of.
              </div>
              <div className="whitespace-pre-wrap bg-surface p-3 rounded-md">
                {appFields.essay2}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  function ReviewerPanel() {
    return (
      <div className="relative p-2 rounded-md border border-input-border bg-card">
        <h3 className="font-medium text-text-secondary mb-3">
          Rubric and Controls
        </h3>
        <RatingFields
          experience={experience}
          passion={passion}
          onExperience={(v) => setExperience(v)}
          onPassion={(v) => setPassion(v)}
        />
        <div className="flex gap-2 mt-4">
          <ReviewNavigation
            mode={mode}
            currentIndex={currentIndex}
            isLast={isLast}
            allowSubmit={allowSubmit}
            onSubmit={handleSubmitReview}
            next={next}
            back={back}
            reset={reset}
          />
        </div>
        <div className="absolute top-0 right-0">
          {autoDecision ? (
            <div className="flex items-center gap-2 p-2">
              <div
                className={`flex items-center gap-2 px-3 py-1 rounded-md border text-sm font-medium ${
                  autoDecision === "auto_accept"
                    ? "bg-emerald-50 border-emerald-200 text-emerald-800"
                    : "bg-rose-50 border-rose-200 text-rose-800"
                }`}
              >
                {autoDecision === "auto_accept" ? <TablerCheck /> : <TablerX />}
                <span>
                  {autoDecision === "auto_accept"
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
            <>
              <div className="flex gap-2 mt-2 ml-2 mr-2">
                <DialogTrigger>
                  <Button className="px-1 h-8" size="sm">
                    <TablerCheck />
                    Auto Accept
                  </Button>

                  <Modal>
                    <JustificationModal accept={true} />
                  </Modal>
                </DialogTrigger>
                <DialogTrigger>
                  <Button className="px-1 h-8" size="sm" variant="danger">
                    <TablerX />
                    Auto Reject
                  </Button>

                  <Modal>
                    <JustificationModal accept={false} />
                  </Modal>
                </DialogTrigger>
              </div>
              {user.role === "staff" && (
                <span className="text-xs text-text-secondary">
                  Requests will be sent to organizers for review
                </span>
              )}
            </>
          )}
        </div>
      </div>
    );
  }

  function JustificationModal({ accept }: { accept: boolean }) {
    const [justification, setJustification] = useState("");

    return (
      <div className="flex flex-col gap-3">
        <label>Justification</label>
        <Input
          placeholder="Type a reason"
          value={justification}
          onChange={(e) => setJustification(e.target.value)}
        />
        <Button
          onClick={() => {
            if (accept) {
              handleRequestAutoAccept(justification);
            } else {
              handleRequestAutoReject(justification);
            }
          }}
          className="px-1 h-8 w-fit"
          size="sm"
          variant={accept ? "primary" : "danger"}
        >
          {accept ? (
            <>
              <TablerCheck /> Auto Accept
            </>
          ) : (
            <>
              <TablerX />
              Auto Reject
            </>
          )}
        </Button>
      </div>
    );
  }

  return (
    <div className="grid lg:grid-cols-2 gap-6 mb-8">
      <div className="space-y-4 mt-4">
        <ApplicantInfo />
        <Essays />
      </div>

      <div className="space-y-4 mt-4">
        <div className="p-2 rounded-md border border-input-border h-[64vh]">
          {resume === "" ? (
            <p>No resume provided.</p>
          ) : (
            <object
              className="w-full h-full"
              type="application/pdf"
              data={resume}
            >
              <p>
                Your browser does not support PDFs.{" "}
                <a href={resume}>Download the PDF</a>.
              </p>
            </object>
          )}
        </div>

        <ReviewerPanel />
      </div>
    </div>
  );
}

interface RatingFieldsProps {
  experience: number;
  passion: number;
  onExperience: (value: number) => void;
  onPassion: (value: number) => void;
}

const RatingFields = ({
  experience,
  passion,
  onExperience,
  onPassion,
}: RatingFieldsProps) => {
  return (
    <div className="w-1/3">
      <div className="flex flex-row gap-4 justify-between items-center">
        <p className="text-lg">Experience:</p>
        <Rating
          style={{ maxWidth: 150, minWidth: 100 }}
          value={experience}
          onChange={onExperience}
        />
      </div>

      <div className="flex flex-row gap-4 items-center justify-between mt-4">
        <p className="text-lg">Passion:</p>
        <Rating
          style={{ maxWidth: 150, minWidth: 100 }}
          value={passion}
          onChange={onPassion}
        />
      </div>
    </div>
  );
};

type ButtonMode = "submit" | "completed-dirty" | "completed-clean";

interface ReviewNavigationProps {
  isLast: boolean;
  mode: ButtonMode;
  allowSubmit: boolean;
  currentIndex: number;
  back: () => void;
  onSubmit: () => void;
  next: () => void;
  reset: () => void;
}

export function ReviewNavigation({
  isLast,
  mode,
  allowSubmit,
  currentIndex,
  back,
  onSubmit,
  next,
  reset,
}: ReviewNavigationProps) {
  const backButton =
    currentIndex > 0 ? (
      <Button
        variant="secondary"
        className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
        onClick={back}
      >
        <TablerArrowLeft />
        Back
      </Button>
    ) : null;

  const buttons = (() => {
    switch (mode) {
      case "submit":
        return (
          <Button
            isDisabled={!allowSubmit}
            className={`flex gap-1 items-center text-lg rounded-md py-2 px-3 ${
              !allowSubmit ? "opacity-30" : ""
            }`}
            variant={isLast ? "success" : "primary"}
            onClick={onSubmit}
          >
            {isLast ? "Finish" : "Submit and Proceed"}
            {isLast ? <TablerCheck /> : <TablerArrowRight />}
          </Button>
        );
      case "completed-clean":
        return (
          <Button
            variant={isLast ? "success" : "secondary"}
            className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
            onClick={next}
          >
            {isLast ? "Finish" : "Next"}
            {isLast ? <TablerCheck /> : <TablerArrowRight />}
          </Button>
        );
      case "completed-dirty":
        return (
          <div className="flex flex-row gap-4">
            <Button
              variant="icon"
              className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
              onClick={reset}
            >
              <TablerRefresh />
            </Button>
            <Button
              variant={isLast ? "success" : "primary"}
              isDisabled={!allowSubmit}
              className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
              onClick={onSubmit}
            >
              {isLast ? "Save and Finish" : "Save and Proceed"}
              {isLast ? <TablerCheck /> : <TablerArrowRight />}
            </Button>
          </div>
        );
    }
  })();

  return (
    <div className="w-full flex flex-row mt-4 gap-3">
      {backButton}
      {buttons}
    </div>
  );
}
