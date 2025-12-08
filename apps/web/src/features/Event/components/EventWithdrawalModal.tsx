import { Modal } from "@/components/ui/Modal";
import { auth } from "@/lib/authClient";
import { Button } from "@/components/ui/Button";
interface EventWithdrawalModalProps {
  event_id: string;
}

function EventWithdrawalModal({ event_id }: EventWithdrawalModalProps) {
  const { data: userData } = auth.useUser();
  const { user } = userData!;
  if (!user) {
    return <div>Loading...{event_id}</div>;
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
          Withdraw Acceptance?
        </h3>
        <p className="mt-1 text-sm text-gray-500">
          Are you sure? This action cannot be undone.
        </p>
        <p className="mt-1 text-sm text-gray-500">
          You can still join the waitlist after withdrawing.
        </p>
      </div>
      <div className="mt-3 flex justify-center gap-3">
        <Button variant="danger">Withdraw</Button>
        {/*When withdrawing, we need to change the accepted -> rejected and then refresh the page*/}
      </div>
    </Modal>
  );
}

export { EventWithdrawalModal };
