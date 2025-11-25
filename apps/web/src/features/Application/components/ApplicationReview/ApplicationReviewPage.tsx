import TablerLoader from "~icons/tabler/loader-2";
import { useAssignedApplication } from "@/features/Application/hooks/useAssignedApplication";
import { Rating } from "@smastrom/react-rating";
import { useState } from "react";
import TablerArrowLeft from "~icons/tabler/arrow-left";
import TablerArrowRight from "~icons/tabler/arrow-right";
import { ProgressBar } from "@/components/ui/ProgressBar";
import { Button } from "@/components/ui/Button";

interface ApplicationReviewPageProps {
  eventId: string;
}

const TODO_appReviewed = true;

// TODO:
// download and display url using presigned urls?
// get all assigned applications (user ids)

export default function ApplicationReviewPage({
  eventId,
}: ApplicationReviewPageProps) {
  const [currentApplicationIdx, setCurrentApplicationIdx] = useState(0);

  const userId = "040d9ebd-9f36-47b9-b330-231cb2e0f095";
  const assignedApplication = useAssignedApplication(eventId, userId);

  const next = () => {
    setCurrentApplicationIdx((prev) => prev + 1);
  };

  const back = () => {
    setCurrentApplicationIdx((prev) => prev - 1);
  };

  if (assignedApplication.error || assignedApplication.isError) {
    return (
      <div className="flex gap-2 flex-col">
        <ProgressBar
          className="w-full"
          label="Review Progress"
          maxValue={100}
          value={currentApplicationIdx + 1}
          valueLabel={`${currentApplicationIdx + 1}/100 applications`}
        />

        <p>An error occurred while loading application for this user.</p>

        <div className="flex gap-2">
          <Button
            variant="secondary"
            className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
            onClick={back}
          >
            <TablerArrowLeft />
            Back
          </Button>
          <Button
            variant="secondary"
            className={`flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3`}
            onClick={next}
          >
            Next
            <TablerArrowRight />
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full gap-4">
      <div className="w-full sm:w-[50%] h-full space-y-5 overflow-y-auto px-2">
        <div>
          <ProgressBar
            className="w-full"
            label="Review Progress"
            maxValue={100}
            value={currentApplicationIdx + 1}
            valueLabel={`${currentApplicationIdx + 1}/100 applications`}
          />
        </div>

        {assignedApplication.isLoading ? (
          <div className="flex gap-2">
            <TablerLoader className="text-xl animate-spin" />
            <p>Loading application</p>
          </div>
        ) : (
          <ApplicationReviewContainer
            next={next}
            back={back}
            eventId={eventId}
            application={JSON.parse(
              atob(assignedApplication.data!["application"] as string),
            )}
          />
        )}
      </div>

      <div className="w-full sm:w-[50%] h-full">
        <iframe
          src="/assets/sample.pdf"
          style={{
            width: "100%",
            height: "100%",
            border: "none",
          }}
        ></iframe>
      </div>
    </div>
  );
}

interface ApplicationReviewContainerProps {
  eventId: string;
  userId?: string;
  application: any;
  next: () => void;
  back: () => void;
}

function ApplicationReviewContainer({
  // eventId,
  // userId = "040d9ebd-9f36-47b9-b330-231cb2e0f095",
  application,
  next,
  back,
}: ApplicationReviewContainerProps) {
  const [passion, setPassion] = useState(0);
  const [experience, setExperience] = useState(0);
  const [hasUpdatedRatings, setHasUpdatedRatings] = useState(false);

  const onSubmitReview = () => {
    if (passion === 0 || experience === 0) {
      alert("Passion or Experience cannot have a 0 rating.");
      return;
    }
  };

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

  return (
    <>
      <div className="p-2 rounded-md border border-input-border">
        <span className="block mb-3">Applicant Information</span>
        <div className="space-y-3">
          <div className="flex gap-8">
            <div className="flex flex-col">
              <span className="text-text-secondary">Name</span>
              <span>
                {application["firstName"] + " " + application["lastName"]}
              </span>
            </div>

            <div className="flex flex-col">
              <span className="text-text-secondary">Major(s)</span>
              <span>
                {Array.isArray(application["majors"])
                  ? application["majors"].join(", ")
                  : application["majors"].split(",").join(", ")}
              </span>
            </div>

            <div className="flex flex-col">
              <span className="text-text-secondary">School</span>
              <span className="truncate max-w-60">{application["school"]}</span>
            </div>

            <div className="flex flex-col">
              <span className="text-text-secondary">Graduation Year</span>
              <span>{application["graduationYear"]}</span>
            </div>
          </div>

          <div className="flex gap-10">
            <div className="flex flex-col">
              <span className="text-text-secondary">
                # of Hackathons Attended
              </span>
              <span>
                {getHackathonExperienceText(application["experience"])}
              </span>
            </div>

            <div className="flex flex-col">
              <span className="text-text-secondary">Project Experience</span>
              <span>
                {getProjectExperienceText(application["projectExperience"])}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-5 p-2 rounded-md border border-input-border">
        <span className="block mb-3">Essay Responses</span>

        <EssayResponse
          question="What is your most memorable experience working in a group? What did
            you learn and accomplish?"
          response={application["essay1"]}
        />

        <EssayResponse
          question=" Tell us about a project you are most proud of."
          response={application["essay2"]}
        />
      </div>

      <div className="flex justify-between">
        <div className="space-y-3">
          <p className="text-lg">Reviewer Ratings</p>

          <div className="flex gap-3 text-lg">
            <div className="flex flex-col w-fit gap-4">
              <span>Experience: </span>
              <span className="self-end">Passion: </span>
            </div>

            <div className="flex gap-7">
              <div className="flex flex-col gap-3">
                <Rating
                  style={{ maxWidth: 150 }}
                  value={experience}
                  onChange={(v: number) => {
                    setExperience(v);
                    setHasUpdatedRatings(true);
                  }}
                />
                <Rating
                  style={{ maxWidth: 150 }}
                  value={passion}
                  onChange={(v: number) => {
                    setPassion(v);
                    setHasUpdatedRatings(true);
                  }}
                />
              </div>
              {TODO_appReviewed && (
                <Button
                  isDisabled={
                    (TODO_appReviewed &&
                      !hasUpdatedRatings &&
                      experience === 0) ||
                    passion === 0
                  }
                  className="mt-3 self-end"
                  onClick={onSubmitReview}
                >
                  Update Review
                </Button>
              )}
            </div>
          </div>
        </div>

        <div className="flex justify-center mt-10 self-end">
          <div className="flex gap-2">
            <Button
              variant="secondary"
              className="flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3"
              onClick={back}
            >
              <TablerArrowLeft />
              Back
            </Button>
            {TODO_appReviewed ? (
              <Button
                variant="secondary"
                className={`flex gap-1 items-center text-lg bg-button-secondary rounded-md py-2 px-3`}
                onClick={next}
              >
                Next
                <TablerArrowRight />
              </Button>
            ) : (
              <Button
                isDisabled={experience === 0 || passion === 0}
                className={`flex gap-1 items-center text-lg rounded-md py-2 px-3 ${(experience === 0 || passion === 0) && "opacity-30"}`}
                onClick={() => {
                  onSubmitReview();
                  next();
                }}
              >
                Submit and Proceed
                <TablerArrowRight />
              </Button>
            )}
          </div>
        </div>
      </div>
    </>
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
