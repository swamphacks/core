import TablerLoader from "~icons/tabler/loader-2";
import { useApplication } from "@/features/Application/hooks/useApplication";
import EssayResponse from "./EssayResponse";
import { RatingFields } from "./RatingFields";
import type { AssignedApplications } from "@/features/Application/hooks/useAssignedApplication";
import { useApplicationResume } from "../../hooks/useAppResume";
import { useRatings } from "../../hooks/useRatings";
import { ReviewNavigation } from "./ReviewNavigation";
import { useAppReviewActions } from "../../hooks/useAppReviewActions";
import { toast } from "react-toastify";

interface ApplicationReviewContainerProps {
  eventId: string;
  assignedApplication: AssignedApplications[number];
  currentIndex: number;
  totalApplications: number;
  next: () => void;
  back: () => void;
}

export default function ApplicationReviewContainer({
  eventId,
  assignedApplication,
  totalApplications,
  currentIndex,
  next,
  back,
}: ApplicationReviewContainerProps) {
  const application = useApplication(eventId, assignedApplication.user_id);
  const resume = useApplicationResume(eventId, assignedApplication.user_id);
  const { review } = useAppReviewActions(eventId, assignedApplication.user_id);
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
  const isLast = currentIndex === totalApplications - 1;
  const isFilled = experience > 0 && passion > 0;
  const allowSubmit = isFilled && isDirty;

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

  const mode = (() => {
    if (!isCompleted) return "submit"; // Not last, not submitted yet
    if (isCompleted && isDirty) return "completed-dirty";
    return "completed-clean"; // Not last, already submitted and clean
  })();

  return (
    <div className="flex flex-row gap-6 mb-8 min-h-[85vh]">
      {/* Application fields */}
      <div className="flex-1 flex flex-col justify-between">
        <div>
          {/* Applicant Information */}
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

        {/* Navigation Buttons */}
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

      {/* PDF Viewer */}
      <div className="flex-1">
        {resume.isLoading ? (
          <p>Resume loading...</p>
        ) : resume.isError || !resume.data ? (
          <p>No resume provided.</p>
        ) : (
          <object
            className="w-full h-full"
            type="application/pdf"
            data={resume.data}
          >
            <p>
              Your browser does not support PDFs.
              <a href={resume.data}>Download the PDF</a>.
            </p>
          </object>
        )}
      </div>
    </div>
  );
}
