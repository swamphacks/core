import { Modal } from "@/components/ui/Modal";
import { auth } from "@/lib/authClient";
import { Button } from "@/components/ui/Button";
interface EventWaitlistModalProps {
  event_id: string;
}

function EventWaitlistModal({ event_id }: EventWaitlistModalProps) {
  const { data: userData } = auth.useUser();
  const { user } = userData!;
  if (!user) {
    return <div>Loading... {event_id}</div>;
  }
  /*     
      {user.userId}
      {event_id}
    */

  return (
    <Modal
      isDismissible
      className="bg-grey p-6 rounded-xl shadow-lg max-w-sm w-full"
    >
      <div className="text-center">
        <h3 className="text-lg font-semibold dark:text-gray-100">
          Join Waitlist?
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          Unfortunately, we weren't able to accomodate everybody who applied
          this year. Join the waitlist, and we will let you know if a spot opens
          up!
        </p>
      </div>
      <div className="mt-3 flex justify-center gap-3">
        <Button variant="success">Join Waitlist</Button>
        {/*When withdrawing, we need to change the accepted -> rejected and then refresh the page*/}
      </div>
    </Modal>
  );
}

export { EventWaitlistModal };
