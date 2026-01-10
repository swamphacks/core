import { Card } from "@/components/ui/Card";
import { generateCheckInIntent } from "@/lib/qr-intents/generate";
import { Heading } from "react-aria-components";
import QRCode from "react-qr-code";

interface Props {
  userId: string;
  eventId: string;
}

export default function AttendeeOverview({ eventId, userId }: Props) {
  const identificationIntentString = generateCheckInIntent(userId, eventId);

  // Used to get the right colors for QR Code
  const styles = getComputedStyle(document.documentElement);
  const bg = styles.getPropertyValue("--surface").trim();
  const fg = styles.getPropertyValue("--text-main").trim();

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Overview
      </Heading>

      <Card className="p-5 w-full md: w-64">
        <QRCode
          bgColor={bg}
          fgColor={fg}
          className="h-full w-full"
          value={identificationIntentString}
        />
      </Card>
    </main>
  );
}
