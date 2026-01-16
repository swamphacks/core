import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import { Scanner } from "@yudiel/react-qr-scanner";
import { useCheckInScanner } from "@/features/CheckIn/hooks/useCheckInScanner";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/check-in",
)({
  component: RouteComponent,
});

function RouteComponent() {
  // const eventId = Route.useParams().eventId;
  const { isScannerActive, onScan } = useCheckInScanner();

  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Check In
      </Heading>

      <Scanner
        classNames={{ container: "w-full max-w-md rounded-lg" }}
        sound={false}
        formats={["qr_code"]}
        allowMultiple
        paused={!isScannerActive}
        onScan={onScan}
        scanDelay={250}
        onError={() => console.log("Something went wrong")}
      />
    </main>
  );
}
