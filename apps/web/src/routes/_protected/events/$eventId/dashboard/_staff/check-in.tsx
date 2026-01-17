import { createFileRoute } from "@tanstack/react-router";
import { Heading } from "react-aria-components";
import { Scanner } from "@yudiel/react-qr-scanner";
import { useCheckInScanner } from "@/features/CheckIn/hooks/useCheckInScanner";
import { useState } from "react";
import { Modal } from "@/components/ui/Modal";
import { Avatar } from "@/components/ui/Avatar";
import { Button } from "@/components/ui/Button";
import { Badge } from "@/components/ui/Badge";

export const Route = createFileRoute(
  "/_protected/events/$eventId/dashboard/_staff/check-in",
)({
  component: RouteComponent,
});


function CheckInBadge({ isCheckedIn }: { isCheckedIn: boolean }) {
  if (isCheckedIn) {
    return (
      <Badge className="bg-green-100 text-green-800 border border-green-300">
        Checked In
      </Badge>
    )
  }

  return (
    <Badge className="bg-yellow-100 text-yellow-800 border border-yellow-300">
      Not Checked In
    </Badge>
  )
}
function RouteComponent() {
  const eventId = Route.useParams().eventId;
  const { isScannerActive, onScan } = useCheckInScanner();
  const [isModalOpen, setIsModalOpen] = useState<boolean>(true);
  const [rfid, setRfid] = useState<string | null>(null);


  return (
    <main>
      <Heading className="text-2xl lg:text-3xl font-semibold mb-6">
        Check In
      </Heading>


      <Modal
        responsive="sheet"
        dialogAriaLabel="checkin-fields"
        isOpen={isModalOpen}
        isDismissible
        onOpenChange={setIsModalOpen}
      >
        <div className="flex flex-col gap-6">

          {/* Header */}
          <div className="flex items-center gap-4">
            <Avatar
              src="https://i.pinimg.com/736x/b1/21/18/b1211897266308c27112a58eacad7482.jpg"
              size="lg"
            />

            <div className="flex flex-col flex-1 min-w-0">
              <h1 className="text-md font-medium truncate">Alexander Wang</h1>
              <p className="text-sm text-text-secondary truncate">
                alex.wg.work@gmail.com
              </p>

              {/* Role + status row */}
              <div className="mt-2 flex flex-wrap gap-2">
                <span className="px-2 py-0.5 rounded-md text-xs font-medium bg-blue-100 text-blue-800 border border-blue-200">
                  Attendee
                </span>

                <CheckInBadge isCheckedIn={false} />
              </div>
            </div>
          </div>


          {/* Mobile warning */}
          <p className="sm:hidden text-sm text-yellow-700 bg-yellow-50 border border-yellow-200 rounded-md px-3 py-2">
            On mobile, check-in requires an attached RFID scanner.
          </p>

          {/* Hidden RFID input */}
          <input
            type="text"
            autoFocus
            onChange={(e) => setRfid(e.target.value)}
            inputMode="none"
            className="opacity-0 h-0 w-0 overflow-hidden"
          />

          { rfid && (<p>{rfid}</p>)}

          {/* Action */}
          <Button
            variant="primary"
            isDisabled={true /* already checked in OR role !== attendee */}
            className="w-full"
          >
            Check In
          </Button>

          {/* Disabled reason */}
          <p className="text-xs text-text-secondary text-center">
            This user cannot be checked in.
          </p>

        </div>
      </Modal>

      <Button onPress={() => setIsModalOpen(true)}>Open Modal </Button>

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
