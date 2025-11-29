import TablerLoader from "~icons/tabler/loader-2";
import { Button } from "@/components/ui/Button";
import { ProgressBar } from "@/components/ui/ProgressBar";
import { useAssignedApplications } from "../../hooks/useAssignedApplications";
import { useAppReviewTutorial } from "../../hooks/useAppReviewTutorial";
import { useAppReviewProgress } from "../../hooks/useAppReviewProgress";
import ApplicationReviewContainer from "./ApplicationReviewContainer";

interface ApplicationReviewPageProps {
  eventId: string;
}

export default function ApplicationReviewPage({
  eventId,
}: ApplicationReviewPageProps) {
  const assignedApps = useAssignedApplications(eventId);
  const appReviewProgress = useAppReviewProgress(assignedApps.data || []);
  const appTutorial = useAppReviewTutorial();

  if (appTutorial.isLoading || assignedApps.isLoading) {
    return (
      <div className="flex gap-2">
        <TablerLoader className="text-xl animate-spin" />
        <p>Loading application review...</p>
      </div>
    );
  }

  if (assignedApps.isError || assignedApps.error) {
    return <div>An error occurred while loading assigned applications.</div>;
  }

  if (!assignedApps.data || assignedApps.data.length === 0) {
    return (
      <div>
        <p className="text-text-secondary">
          You have no assigned applications to review at this time. Refresh the
          page?
          <br />
          Ask your organizer if you believe this is an error.
        </p>
      </div>
    );
  }

  if (!appTutorial.isCompleted) {
    return (
      <div className="w-full pb-8">
        <p className="mb-8 text-text-secondary">
          Welcome to the application review process! Follow these steps to get
          started:
        </p>

        {/* Step cards container */}
        <div className="flex flex-col md:flex-row gap-6 w-full">
          {/* Step 1 */}
          <div className="flex-1 bg-surface p-6 rounded-lg shadow-sm flex flex-col items-start">
            <img
              src="https://static.vecteezy.com/system/resources/previews/049/674/041/non_2x/a-scroll-with-a-feather-and-a-quill-free-png.png"
              alt="Step 1"
              className="w-full h-32 object-contain mb-4"
            />
            <h3 className="text-xl font-semibold mb-2">
              Step 1: Read Carefully
            </h3>
            <p className="text-text-secondary">
              Carefully read each application before rating. The application
              fields will be on the left side and their resume will be on the
              right. The rating star buttons will be underneath the application
              fields on the left.
            </p>
          </div>

          {/* Step 2 */}
          <div className="flex-1 bg-surface p-6 rounded-lg shadow-sm flex flex-col items-start">
            <img
              src="https://pngimg.com/d/star_PNG76902.png"
              alt="Step 2"
              className="w-full h-32 object-contain mb-4"
            />
            <h3 className="text-xl font-semibold mb-2">
              Step 2: Rate Applicants
            </h3>
            <p className="text-text-secondary">
              Applicants should be rated on two separate scales: Experience and
              Passion. Use the star ratings to provide your evaluation.
              Experience reflects the applicant&apos;s relevant skills and
              background, while Passion measures their enthusiasm and commitment
              to the field.
            </p>
          </div>

          {/* Step 3 */}
          <div className="flex-1 bg-surface p-6 rounded-lg shadow-sm flex flex-col items-start">
            <img
              src="https://www.freeiconspng.com/uploads/save-icon--1.png"
              alt="Step 3"
              className="w-full h-32 object-contain mb-4"
            />
            <h3 className="text-xl font-semibold mb-2">Step 3: Navigate</h3>
            <p className="text-text-secondary">
              Use the "Submit/Next" and "Back" buttons to navigate between
              applications. Make sure to complete reviews for all assigned
              applications. You may go back and adjust your ratings, but be sure
              to save your changes!
            </p>
          </div>
        </div>

        {/* Complete tutorial button */}
        <div className="mt-8">
          <Button onClick={appTutorial.completeTutorial}>
            Complete Tutorial
          </Button>
        </div>
      </div>
    );
  }

  if (appReviewProgress.finished) {
    return (
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
          You have completed reviewing all assigned applications. Thank you for
          your time and effort! You can go back to change your reviews if
          needed.
        </p>

        <div className="flex flex-row gap-4">
          <Button variant="secondary" onClick={appReviewProgress.goPrevious}>
            Change Reviews
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex h-full gap-4">
      <div className="w-full hidden md:flex md:flex-col">
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

      <div className="flex md:hidden">
        <p className="text-lg text-red-600">
          Please switch to a laptop or desktop to continue.
        </p>
      </div>
    </div>
  );
}
