import { Button } from "@/components/ui/Button";
import TablerFolder from "~icons/tabler/folder-open";
import TablerScan from "~icons/tabler/scan";
import { type ColumnDef, type FilterFn } from "@tanstack/react-table";
import { useMemo, useState } from "react";
import { DialogTrigger, TooltipTrigger, Tooltip } from "react-aria-components";

import { Route as EventUsersRoute } from "@/routes/_protected/events/$eventId/dashboard/_staff/attendee-directory";

import { Table } from "@/components/ui/Table";
import type { EventUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import { useUrlTableState } from "@/features/EventAdmin/hooks/useUrlTableState";
import RoleBadge from "@/features/EventAdmin/components/RoleBadge";
import { UserSideDrawer } from "@/features/EventAdmin/components/UserSideDrawer";
import CheckInRFIDModal from "@/features/CheckIn/components/CheckInRFIDModal";

// Warning: When using the url table state saving (useUrlTableState hook), random query parameters may be interpreted as table filters if column name is identical.

// Filtering for event user data type
const fuzzyTextFilterFn: FilterFn<EventUser> = (row, columnId, value) => {
  const rowValue = row.getValue(columnId) as string;
  return rowValue.toLowerCase().includes((value as string).toLowerCase());
};

// Filter function for checked-in status (yes/no)
const checkedInFilterFn: FilterFn<EventUser> = (row, _columnId, value) => {
  const checkedInAt = row.original.checked_in_at;
  const isCheckedIn = checkedInAt !== null && checkedInAt !== undefined;
  const filterValues = value as string[];

  if (!filterValues || filterValues.length === 0) return true;

  return filterValues.some((val) => {
    if (val === "yes") return isCheckedIn;
    if (val === "no") return !isCheckedIn;
    return false;
  });
};

const fallbackData: EventUser[] = [];

interface Props {
  eventId: string;
  data?: EventUser[];
}

const AttendeeTable = ({ data, eventId }: Props) => {
  // For use with url filter state
  const search = EventUsersRoute.useSearch();
  const navigate = EventUsersRoute.useNavigate();

  const [checkInModalOpen, setCheckInModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<EventUser | null>(null);

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

  const handleOpenCheckIn = (user: EventUser) => {
    setSelectedUser(user);
    setCheckInModalOpen(true);
  };

  const handleCloseCheckIn = () => {
    setCheckInModalOpen(false);
    setSelectedUser(null);
  };

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
        meta: {
          responsiveClass: "hidden xl:table-cell",
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
        accessorFn: (row) => row.preferred_email ?? row.email,
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
        id: "checked_in",
        header: "Checked In",
        accessorFn: (row) => {
          return row.checked_in_at !== null && row.checked_in_at !== undefined
            ? "Yes"
            : "No";
        },
        cell: ({ row }) => {
          const checkedInAt = row.original.checked_in_at;
          const isCheckedIn = checkedInAt !== null && checkedInAt !== undefined;
          return (
            <span
              className={
                isCheckedIn
                  ? "text-green-600 dark:text-green-400"
                  : "text-gray-500"
              }
            >
              {isCheckedIn ? "Yes" : "No"}
            </span>
          );
        },
        filterFn: checkedInFilterFn,
        enableSorting: true,
        sortingFn: "alphanumeric",
        meta: {
          filterType: "select",
          filterOptions: [
            { value: "yes", label: "Yes" },
            { value: "no", label: "No" },
          ],
        },
      },
      {
        id: "actions",
        header: "",
        size: 50,
        cell: ({ row }) => {
          return (
            <div className="flex items-center gap-1.5 justify-end pr-2">
              <DialogTrigger>
                <TooltipTrigger delay={250} closeDelay={250}>
                  <Button
                    variant="primary"
                    className="aspect-square p-2 shrink-0"
                  >
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
              <TooltipTrigger delay={250} closeDelay={250}>
                <Button
                  variant="secondary"
                  className="aspect-square p-2 shrink-0"
                  onPress={() => handleOpenCheckIn(row.original)}
                >
                  <TablerScan className="h-4 w-4" />
                </Button>
                <Tooltip
                  offset={5}
                  className="bg-surface border-input-border border-2 flex justify-center items-center py-1 px-2 rounded-md"
                >
                  Check In with RFID
                </Tooltip>
              </TooltipTrigger>
            </div>
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
        key="Attendee Management"
        data={data ?? fallbackData}
        columns={columns}
        columnFilters={columnFilters}
        onColumnFiltersChange={setColumnFilters}
        sorting={sorting}
        onSortingChange={setSorting}
        pagination={pagination}
        onPaginationChange={setPagination}
      />

      {selectedUser && (
        <CheckInRFIDModal
          isOpen={checkInModalOpen}
          onClose={handleCloseCheckIn}
          user={selectedUser}
          eventId={eventId}
        />
      )}
    </div>
  );
};

export default AttendeeTable;
