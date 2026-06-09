import { useApplicationStatistics } from "@/modules/Application/hooks/useApplicationStatistics";
import type { StaffHackathon } from "@/lib/openapi/types";
import { Button } from "@/components/ui/Button";

interface ReviewNotStartedProps {
  hackathon: StaffHackathon;
  startApplicationReview: () => void;
}

export default function ReviewNotStartedAdmin({
  hackathon,
  startApplicationReview,
}: ReviewNotStartedProps) {
  const stats = useApplicationStatistics();

  if (stats.isLoading || !stats.data) {
    return <div>Loading...</div>;
  }

  const now = new Date();
  const applicationPeriodClosed = now >= new Date(hackathon.application_close);
  const validNumOfApplicants = stats.data.statusStats.submitted > 0;

  return (
    <div className="flex flex-col gap-4">
      <p className="text-text-secondary">Application review has not started.</p>

      <Button
        variant="primary"
        className="w-fit"
        isDisabled={!(validNumOfApplicants && applicationPeriodClosed)}
        onClick={startApplicationReview}
      >
        Start Application Review
      </Button>

      {!applicationPeriodClosed ? (
        <p className="text-text-secondary text-sm">
          The application period is still ongoing.
        </p>
      ) : (
        !validNumOfApplicants && (
          <p className="text-text-secondary text-sm">
            There are no submitted applications to review yet.
          </p>
        )
      )}
    </div>
  );
}
