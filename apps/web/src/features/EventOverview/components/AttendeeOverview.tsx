import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EventAttendanceWithdrawalModal } from "@/features/Event/components/EventAttendanceWithdrawalModal";
import { generateIdentifyIntent } from "@/lib/qr-intents/generate";
import { DialogTrigger, Heading } from "react-aria-components";
import QRCode from "react-qr-code";

interface Props {
  userId: string;
  eventId: string;
}

export default function AttendeeOverview({ userId, eventId }: Props) {
  const identificationIntentString = generateIdentifyIntent(userId);

  // Used to get the right colors for QR Code
  const styles = getComputedStyle(document.documentElement);
  const bg = styles.getPropertyValue("--surface").trim();
  const fg = styles.getPropertyValue("--text-main").trim();

  // Should be a feature flag
  const allowWithdrawal = false;

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Overview
      </Heading>
      <Card className="p-5 w-full md:w-64">
        <QRCode
          bgColor={bg}
          fgColor={fg}
          className="h-full w-full"
          value={identificationIntentString}
        />
      </Card>
      {allowWithdrawal && (
        <div>
          <p className="my-5">Can't make it to the event?</p>
          <DialogTrigger>
            <Button variant="danger">{"Withdraw Attendance"}</Button>
            <EventAttendanceWithdrawalModal
              eventId={eventId}
            ></EventAttendanceWithdrawalModal>
          </DialogTrigger>
        </div>
      )}
    </main>
  );
}
