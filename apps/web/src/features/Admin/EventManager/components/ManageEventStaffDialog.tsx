import { Heading } from "react-aria-components";
import { useEventStaffUsers } from "../hooks/useEventStaffUsers";
import { useAdminStaffActions } from "../hooks/useAdminStaffActions";
import { AddStaffForm } from "./AddStaffForm";
import { StaffTable } from "./StaffTable";
import { type ColumnDef } from "@tanstack/react-table";
import { type StaffUser } from "../hooks/useEventStaffUsers";
import { Modal } from "@/components/ui/Modal";

const columns: ColumnDef<StaffUser>[] = [
  {
    id: "name",
    header: "Name",
    cell: ({ row }) => {
      const { name } = row.original;
      return name?.trim() || "Unknown";
    },
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "eventRole",
    header: "Role",
  },
];

export function ManageEventStaffDialog({ eventId }: { eventId: string }) {
  const {
    data: staffUsers = [],
    isLoading,
    isError,
  } = useEventStaffUsers(eventId);
  const { add } = useAdminStaffActions(eventId);

  return (
    <Modal isDismissible size="xl">
      <div className="flex flex-col gap-4">
        <Heading className="text-lg" slot="title">
          Manage Users
        </Heading>

        <StaffTable
          data={staffUsers}
          columns={columns}
          isLoading={isLoading}
          isError={isError}
        />

        <AddStaffForm onSubmit={(data) => add.mutateAsync(data)} />
      </div>
    </Modal>
  );
}
