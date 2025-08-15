import { Heading } from "react-aria-components";
import { useEventStaffUsers } from "../hooks/useEventStaffUsers";
import {
  useAdminStaffActions,
  type AssignStaffRole,
} from "../hooks/useAdminStaffActions";
import { AddStaffForm } from "./AddStaffForm";
import { StaffTable } from "./StaffTable";
import { type ColumnDef } from "@tanstack/react-table";
import { Modal } from "@/components/ui/Modal";
import type { User } from "@/lib/openapi/types";
import { HTTPError } from "ky";
import { showToast } from "@/lib/toast/toast";

const columns: ColumnDef<User>[] = [
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
    accessorKey: "event_role",
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

  const addNewStaff = async (data: AssignStaffRole) => {
    await add.mutateAsync(data, {
      onError: (error) => {
        if (error instanceof HTTPError && error.response.status === 404) {
          showToast({
            title: "Error",
            message: `Couldn't find user with email ${data.email}`,
            type: "error",
          });
        }
      },
    });
  };

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

        <AddStaffForm onSubmit={addNewStaff} />
      </div>
    </Modal>
  );
}
