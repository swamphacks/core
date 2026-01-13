import { DialogTrigger } from "react-aria-components";
import { EventAttendanceWithdrawalModal } from "@/features/Event/components/EventAttendanceWithdrawalModal";
import { Button } from "@/components/ui/Button";

interface ApplicationOverviewProps {
  eventId: string;
}

export default function AttendeeOverview({
  eventId,
}: ApplicationOverviewProps) {
  return (
    <div>
      <div>More here coming soon!</div>
      <div>
        <p className="my-5">Can't make it to the event?</p>
        <DialogTrigger>
          <Button variant="danger">{"Withdraw Attendance"}</Button>
          <EventAttendanceWithdrawalModal
            eventId={eventId}
          ></EventAttendanceWithdrawalModal>
        </DialogTrigger>
      </div>
    </div>
  );
}
