import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import {
  Group,
  Heading,
  OverlayTriggerStateContext,
  Text,
} from "react-aria-components";
import { useContext } from "react";
import type { StaffUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";
import { useStaffActions } from "../hooks/useStaffActions";
import { toast } from "react-toastify";

interface DeleteStaffDialogProps {
  eventId: string;
  user: StaffUser;
}

function DeleteStaffDialog({ eventId, user }: DeleteStaffDialogProps) {
  const state = useContext(OverlayTriggerStateContext)!;
  const {
    remove: { mutateAsync, isPending },
  } = useStaffActions(eventId);

  const handleDelete = async () => {
    await mutateAsync(user.id, {
      onSuccess: () => {
        toast.success("Staff member removed successfully.");
        state.close();
      },
      onError: () => {
        toast.error("Failed to remove staff member. Please try again.");
      },
    });
  };

  return (
    <Modal isDismissible>
      <div className="flex flex-col gap-8 justify-center items-center">
        <div className="flex flex-col gap-4">
          <Heading className="text-lg" slot="title">
            Remove {user.name} from event staff?
          </Heading>
          <Text className="text-base">
            Are you sure you want to remove this user from the event staff? You
            can always add them back later.
          </Text>
        </div>
        <div className="w-full flex flex-row justify-end">
          <Group className="gap-4 flex flex-row">
            <Button isPending={isPending} variant="secondary" slot="close">
              Nevermind
            </Button>

            <Button
              onClick={handleDelete}
              isPending={isPending}
              variant="danger"
            >
              Remove
            </Button>
          </Group>
        </div>
      </div>
    </Modal>
  );
}

export { DeleteStaffDialog };
