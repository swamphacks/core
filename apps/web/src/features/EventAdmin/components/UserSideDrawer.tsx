import { Modal } from "@/components/ui/Modal";
import type { StaffUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

interface UserSideDrawerProps {
  user: StaffUser;
  event_id: string;
}
function UserSideDrawer({ user, event_id }: UserSideDrawerProps) {
  return (
    <Modal isDismissible>
      {user.name}
      {event_id}
    </Modal>
  );
}

export { UserSideDrawer };
