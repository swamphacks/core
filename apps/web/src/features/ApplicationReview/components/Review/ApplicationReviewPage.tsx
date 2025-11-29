import TablerLoader from "~icons/tabler/loader-2";
import { useEffect } from "react";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import TablerArrowRight from "~icons/tabler/arrow-right";
import { ProgressBar } from "@/components/ui/ProgressBar";
import { Button } from "@/components/ui/Button";
import {
  useAssignedApplications,
  type AssignedApplications,
} from "../../hooks/useAssignedApplications";
import { useAppReviewTutorial } from "../../hooks/useAppReviewTutorial";
import { useAppReviewProgress } from "../../hooks/useAppReviewProgress";
import { useApplication } from "@/features/Application/hooks/useApplication";
import { Rating } from "@smastrom/react-rating";
import { useRatings } from "../../hooks/useRatings";

interface ApplicationReviewPageProps {
  eventId: string;
}

export default function ApplicationReviewPage({
  eventId,
}: ApplicationReviewPageProps) {
  const assignedApps = useAssignedApplications(eventId);
  const appReviewProgress = useAppReviewProgress(assignedApps.data || []);
  const appTutorial = useAppReviewTutorial();

  useEffect(() => {
    console.log(
      "Current Assigned Application:",
      appReviewProgress.currentAssignedApplication,
    );
    console.log("Current Index:", appReviewProgress.currentIndex);
  }, [
    appReviewProgress.currentAssignedApplication,
    appReviewProgress.currentIndex,
  ]);

  if (appTutorial.isLoading || assignedApps.isLoading) {
    return (
      <div className="flex gap-2">
        <TablerLoader className="text-xl animate-spin" />
        <p>Loading application review...</p>
      </div>
    );
  }

  if (assignedApps.isError || assignedApps.error || !assignedApps.data) {
    return <div>An error occurred while loading assigned applications.</div>;
  }

  if (!appTutorial.isCompleted) {
    return (
      <div>
        <h2 className="text-2xl font-semibold mb-4">
          Application Review Tutorial
        </h2>
        <p className="mb-4">
          Welcome to the application review process! Here are some guidelines to
          help you get started:
        </p>
        <ul className="list-disc list-inside mb-4">
          <li>Carefully read each application before making a decision.</li>
          <li>Rate applicants based on their passion and experience.</li>
          <li>Use the navigation buttons to move between applications.</li>
        </ul>
        <Button onClick={appTutorial.completeTutorial}>
          Complete Tutorial
        </Button>
      </div>
    );
  }

  if (appReviewProgress.finished) {
    return (
      <div>You have completed all assigned application reviews. Thank you!</div>
    );
  }

  return (
    <div className="flex h-full gap-4">
      <div className="w-full">
        <div className="w-[50%] pr-3">
          <ProgressBar
            className="w-full"
            label="Review Progress"
            maxValue={appReviewProgress.totalApplications}
            value={appReviewProgress.currentIndex + 1}
            valueLabel={`${appReviewProgress.currentIndex + 1}/${appReviewProgress.totalApplications} applications`}
          />
        </div>
        {appReviewProgress.currentAssignedApplication ? (
          <ApplicationReviewContainer
            next={appReviewProgress.goNext}
            back={appReviewProgress.goPrevious}
            currentIndex={appReviewProgress.currentIndex}
            totalApplications={appReviewProgress.totalApplications}
            eventId={eventId}
            assignedApplication={appReviewProgress.currentAssignedApplication}
          />
        ) : (
          <div>Loading application...</div>
        )}
      </div>
    </div>
  );
}

interface ApplicationReviewContainerProps {
  eventId: string;
  userId?: string;
  assignedApplication: AssignedApplications[number];
  currentIndex: number;
  totalApplications: number;
  next: () => void;
  back: () => void;
}

function ApplicationReviewContainer({
  eventId,
  assignedApplication,
  currentIndex,
  next,
  back,
}: ApplicationReviewContainerProps) {
  const application = useApplication(eventId, assignedApplication.user_id);
  const { experience, passion, isDirty, setExperience, setPassion, reset } =
    useRatings(
      application.data?.experience_rating || 0,
      application.data?.passion_rating || 0,
    );

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
    }
  };

  if (application.isLoading) {
    return (
      <div className="flex gap-2">
        <TablerLoader className="text-xl animate-spin" />
        <p>Loading application...</p>
      </div>
    );
  }

  if (application.isError || application.error || !application.data) {
    return <div>An error occurred while loading the application.</div>;
  }

  const appFields = application.data.application;
  const isCompleted = assignedApplication.status === "completed";
  const allowSubmit = experience > 0 && passion > 0 && isDirty;

  const mode = (() => {
    if (!isCompleted) {
      return "submit";
    }
    if (isCompleted && isDirty) {
      return "completed-dirty";
    }
    return "completed-clean";
  })();

  return (
    <div className="flex flex-row gap-6 mb-8 min-h-[85vh]">
      {/* Application fields */}
      <div className="flex-1 flex flex-col justify-between">
        {/* Applicant Information */}
        <div>
          <div className="p-2 rounded-md border border-input-border mt-4 mb-4">
            <span className="block mb-3">Applicant Information</span>
            <div className="space-y-3">
              <div className="flex gap-8">
                <div className="flex flex-col">
                  <span className="text-text-secondary">Name</span>
                  <span>{appFields.firstName + " " + appFields.lastName}</span>
                </div>

                <div className="flex flex-col">
                  <span className="text-text-secondary">Major(s)</span>
                  <span>{appFields.majors}</span>
                </div>

                <div className="flex flex-col">
                  <span className="text-text-secondary">School</span>
                  <span className="truncate max-w-60">{appFields.school}</span>
                </div>

                <div className="flex flex-col">
                  <span className="text-text-secondary">Graduation Year</span>
                  <span>{appFields.graduationYear}</span>
                </div>
              </div>

              <div className="flex gap-10">
                <div className="flex flex-col">
                  <span className="text-text-secondary">
                    # of Hackathons Attended
                  </span>
                  <span>
                    {getHackathonExperienceText(appFields.experience)}
                  </span>
                </div>

                <div className="flex flex-col">
                  <span className="text-text-secondary">
                    Project Experience
                  </span>
                  <span>
                    {getProjectExperienceText(appFields.projectExperience)}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Essay Responses */}

          <div className="space-y-5 p-2 rounded-md border border-input-border mb-4">
            <span className="block mb-3">Essay Responses</span>

            <EssayResponse
              question="What is your most memorable experience working in a group? What did
              you learn and accomplish?"
              response={appFields.essay1}
            />

            <EssayResponse
              question=" Tell us about a project you are most proud of."
              response={appFields.essay2}
            />
          </div>

          {/* Reviewer Ratings */}

          <RatingFields
            experience={experience}
            passion={passion}
            onExperience={(v) => setExperience(v)}
            onPassion={(v) => setPassion(v)}
          />
        </div>
        <div className="w-full flex flex-row justify-between mt-4">
          {currentIndex > 0 ? (
            <Button
              variant="secondary"
              className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
              onClick={back}
            >
              <TablerArrowLeft />
              Back
            </Button>
          ) : (
            <div />
          )}

          {mode === "submit" && (
            <Button
              isDisabled={!allowSubmit}
              className={`flex gap-1 items-center text-lg rounded-md py-2 px-3 ${
                !allowSubmit && "opacity-30"
              }`}
              onClick={next}
            >
              Submit and Proceed
              <TablerArrowRight />
            </Button>
          )}

          {mode === "completed-clean" && (
            <Button
              variant="secondary"
              className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
              onClick={next}
            >
              Next
              <TablerArrowRight />
            </Button>
          )}

          {mode === "completed-dirty" && (
            <div className="flex flex-row gap-4">
              <Button
                className="flex gap-1 items-center text-lg rounded-md py-2 px-3"
                onClick={() => {
                  reset();
                }}
              >
                Save Changes
              </Button>

              <Button
                variant="secondary"
                className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
                onClick={() => {
                  // Reset ratings on navigation
                  reset();
                  next();
                }}
              >
                Next
                <TablerArrowRight />
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* PDF Viewer */}
      <div className="flex-1">
        <object
          className="w-full h-full"
          type="application/pdf"
          data="/assets/sample.pdf"
        >
          <p>
            Your browser does not support PDFs.
            <a href="/assets/sample.pdf">Download the PDF</a>.
          </p>
        </object>
      </div>
    </div>
  );
}

interface EssayResponse {
  question: string;
  response: string;
}

function EssayResponse({ question, response }: EssayResponse) {
  return (
    <div>
      <p className="text-text-secondary">{question}</p>
      <p>{response}</p>
    </div>
  );
}

interface RatingFieldsProps {
  experience: number;
  passion: number;
  onExperience: (value: number) => void;
  onPassion: (value: number) => void;
}

export const RatingFields = ({
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
          style={{ maxWidth: 150 }}
          value={experience}
          onChange={onExperience}
        />
      </div>

      <div className="flex flex-row gap-4 items-center justify-between mt-4">
        <p className="text-lg">Passion:</p>
        <Rating
          style={{ maxWidth: 150 }}
          value={passion}
          onChange={onPassion}
        />
      </div>
    </div>
  );
};
