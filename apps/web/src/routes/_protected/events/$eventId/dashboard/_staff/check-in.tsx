import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import { Scanner, type IDetectedBarcode } from "@yudiel/react-qr-scanner";
import { useState } from "react";
import CheckInModal from "@/features/CheckIn/components/CheckInModal";
import { parseQrIntent } from "@/lib/qr-intents/parse";
import { Intent } from "@/lib/qr-intents/intent";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/check-in",
)({
  component: RouteComponent,
});

function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const [scannedUserId, setScannedUserId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState<boolean>(false);

  const onScan = (scannedData: IDetectedBarcode[]) => {
    const res = parseQrIntent(scannedData[0].rawValue);
    if (res.ok && res.value.intent === Intent.IDENT) {
      setScannedUserId(res.value.user_id);
      setIsModalOpen(true);
    }
  };

  return (
    <main className="flex flex-col">
      <header>
        <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
          Check In
        </Heading>
      </header>

      {/* Scanner: Fixed width, no extra fluff */}
      <section className="w-full">
        <div className="overflow-hidden rounded-md bg-black aspect-square w-full sm:w-1/2 lg:w-1/4">
          <Scanner
            classNames={{ container: "w-full h-full" }}
            sound={false}
            formats={["qr_code"]}
            paused={isModalOpen}
            onScan={onScan}
            scanDelay={500}
            onError={(err) => console.log(err)}
          />
        </div>
        <p className="text-sm text-text-secondary my-6">
          Point the camera at a user&apos;s QR code to begin check-in.
        </p>
      </section>

      {/* Future scans table */}
      {/* <section className="w-full pt-8 border-t border-border/40">
        <Heading className="text-lg font-medium mb-4">Recent Activity</Heading>
        <div className="text-sm text-muted-foreground">
          No recent scans to display.
        </div>
      </section> */}

      <CheckInModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        userId={scannedUserId}
        eventId={eventId}
      />
    </main>
  );
}
