import { Modal } from "@/components/ui/Modal";
import type { StaffUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

interface UserSideDrawerProps {
  user: StaffUser;
}
function UserSideDrawer({ user }: UserSideDrawerProps) {
  return <Modal isDismissible>{user.name}</Modal>;
}

export { UserSideDrawer };
