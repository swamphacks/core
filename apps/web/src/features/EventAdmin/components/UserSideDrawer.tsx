import { Modal } from "@/components/ui/Modal";
import type { EventUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";

interface UserSideDrawerProps {
  user: EventUser;
  event_id: string;
}
function UserSideDrawer({ user, event_id }: UserSideDrawerProps) {
  return (
    <Modal isDismissible>
      Under Construction...
      {user.name}
      {event_id}
    </Modal>
  );
}

export { UserSideDrawer };
