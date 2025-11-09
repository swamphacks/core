import { Button } from "@/components/ui/Button";
import TablerFolder from "~icons/tabler/folder-open";
import { type ColumnDef, type FilterFn } from "@tanstack/react-table";
import { useMemo } from "react";
import { DialogTrigger, TooltipTrigger, Tooltip } from "react-aria-components";

import RoleBadge from "./RoleBadge";
import type { EventUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import { UserSideDrawer } from "./UserSideDrawer";
import { Route as EventUsersRoute } from "@/routes/_protected/events/$eventId/dashboard/_admin/user-management";

import { Table } from "@/components/ui/Table";
import { useUrlTableState } from "../hooks/useUrlTableState";

// Filtering for event user data type
const fuzzyTextFilterFn: FilterFn<EventUser> = (row, columnId, value) => {
  const rowValue = row.getValue(columnId) as string;
  return rowValue.toLowerCase().includes((value as string).toLowerCase());
};

const fallbackData: EventUser[] = [];

interface Props {
  eventId: string;
  data?: EventUser[];
}

const UserTable = ({ data, eventId }: Props) => {
  // For use with url filter state
  const search = EventUsersRoute.useSearch();
  const navigate = EventUsersRoute.useNavigate();

  const {
    columnFilters,
    setColumnFilters,
    sorting,
    setSorting,
    pagination,
    setPagination,
  } = useUrlTableState({
    search,
    navigate,
  });

  const columns: ColumnDef<EventUser>[] = useMemo(
    () => [
      {
        id: "avatar",
        header: "Avatar",
        size: 30,
        cell: ({ row }) => {
          const avatarUrl = row.original.image;
          return avatarUrl ? (
            <img
              src={avatarUrl}
              alt="Avatar"
              className="w-full h-auto max-w-[40px] rounded-full object-cover"
            />
          ) : (
            <div className="h-10 w-10 rounded-full bg-gray-300 dark:bg-neutral-700 flex items-center justify-center">
              <span className="text-gray-600 dark:text-neutral-400">N/A</span>
            </div>
          );
        },
      },
      {
        id: "name",
        header: "Name",
        cell: ({ row }) => {
          const name = row.original.name;
          return name ? name : "Unknown";
        },
        accessorKey: "name",
        filterFn: fuzzyTextFilterFn,
        sortingFn: "alphanumeric",
        meta: {
          filterType: "text",
        },
      },
      {
        accessorKey: "email",
        header: "Email",
        filterFn: fuzzyTextFilterFn,
        sortingFn: "alphanumeric",
        meta: {
          filterType: "text",
        },
      },
      {
        accessorKey: "event_role",
        header: "Role",
        cell: ({ row }) => {
          const role = row.original.event_role;
          if (!role) return "N/A";

          return <RoleBadge role={role} />;
        },
        filterFn: "arrIncludesSome",
        enableSorting: false,
        meta: {
          filterType: "select",
          filterOptions: [
            { value: "admin", label: "Admin" },
            { value: "staff", label: "Staff" },
            { value: "attendee", label: "Attendee" },
            { value: "applicant", label: "Applicant" },
          ],
        },
      },
      {
        id: "open-details",
        header: "",
        size: 20,
        cell: ({ row }) => {
          return (
            <DialogTrigger>
              <TooltipTrigger delay={250} closeDelay={250}>
                <Button variant="primary" className="aspect-square p-2">
                  <TablerFolder className="h-4 w-4" />
                </Button>
                <Tooltip
                  offset={5}
                  className="bg-surface border-input-border border-2 flex justify-center items-center py-1 px-2 rounded-md"
                >
                  Open User Details
                </Tooltip>
              </TooltipTrigger>
              <UserSideDrawer user={row.original} event_id={eventId} />
            </DialogTrigger>
          );
        },
      },
    ],
    [eventId],
  );

  return (
    <div>
      <div className="flex justify-end pr-4">
        <Button>Actions</Button>
      </div>

      <Table
        data={data ?? fallbackData}
        columns={columns}
        columnFilters={columnFilters}
        onColumnFiltersChange={setColumnFilters}
        sorting={sorting}
        onSortingChange={setSorting}
        pagination={pagination}
        onPaginationChange={setPagination}
      />
    </div>
  );
};

export default UserTable;
