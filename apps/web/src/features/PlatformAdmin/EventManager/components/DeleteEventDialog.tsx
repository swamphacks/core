import { Button } from "@/components/ui/Button";
import { Modal } from "@/components/ui/Modal";
import type { Event } from "@/lib/openapi/types";
import {
  Group,
  Heading,
  OverlayTriggerStateContext,
  Text,
} from "react-aria-components";
import { useAdminEventActions } from "../hooks/useAdminEventActions";
import { useContext } from "react";

interface DeleteEventDialogProps {
  event: Event;
}

function DeleteEventDialog({ event }: DeleteEventDialogProps) {
  const state = useContext(OverlayTriggerStateContext)!;
  const {
    remove: { mutateAsync, isPending },
  } = useAdminEventActions();
  const { name, id } = event;

  const deleteEvent = async () => {
    await mutateAsync(id, {
      // TODO: This state.close doesn't animate the fadeout correctly.
      onSuccess: () => state.close(),
    });
  };

  return (
    <Modal isDismissible>
      <div className="flex flex-col gap-8 justify-center items-center">
        <div className="flex flex-col gap-4">
          <Heading className="text-lg" slot="title">
            Delete {name}
          </Heading>
          <Text className="text-base">
            Are you sure you want to delete this event? This action is
            irreversable.
          </Text>
        </div>
        <div className="w-full flex flex-row justify-end">
          <Group className="gap-4 flex flex-row">
            <Button isPending={isPending} variant="secondary" slot="close">
              Nevermind
            </Button>

            <Button
              onClick={deleteEvent}
              isPending={isPending}
              variant="danger"
            >
              Delete
            </Button>
          </Group>
        </div>
      </div>
    </Modal>
  );
}

export { DeleteEventDialog };
