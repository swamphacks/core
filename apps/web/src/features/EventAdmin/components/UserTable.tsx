import { Button } from "@/components/ui/Button";
import TablerFolder from "~icons/tabler/folder-open";

import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  useReactTable,
  type ColumnDef,
  type FilterFn,
} from "@tanstack/react-table";
import RoleBadge from "./RoleBadge";
import type { EventUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventUsers";
import { TextField } from "@/components/ui/TextField";
import { useMemo, useState, useEffect } from "react";
import { DialogTrigger, TooltipTrigger, Tooltip } from "react-aria-components";
import { UserSideDrawer } from "./UserSideDrawer";
import debounce from "lodash.debounce";
import { Route as EventUsersRoute } from "@/routes/_protected/events/$eventId/dashboard/_admin/user-management";
const fuzzyTextFilterFn: FilterFn<EventUser> = (row, columnId, value) => {
  const rowValue = row.getValue(columnId) as string;
  return rowValue.toLowerCase().includes((value as string).toLowerCase());
};

interface ColumnFilter {
  id: string;
  value: unknown;
}
type ColumnFiltersState = ColumnFilter[];

const fallbackData: EventUser[] = [];

interface Props {
  eventId: string;
  data?: EventUser[];
}

// TODO: move to utils
function parseEncodedFilters(
  encodedString: string | undefined,
): ColumnFiltersState {
  if (encodedString) {
    try {
      const decoded = atob(encodedString);
      return JSON.parse(decoded);
    } catch (e) {
      /* fall through */
      console.error("Failed to decode filters from URL:", e);
    }
  }
  return [];
}

// TODO: move to utils
function encodeFiltersForUrl(filters: ColumnFiltersState): string {
  return btoa(JSON.stringify(filters));
}

const UserTable = ({ data, eventId }: Props) => {
  const searchParams = EventUsersRoute.useSearch();
  const navigate = EventUsersRoute.useNavigate();

  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>(() => {
    return parseEncodedFilters(searchParams.filters);
  });

  const debouncedUrlUpdate = useMemo(
    () =>
      debounce((filters: ColumnFiltersState, search: typeof searchParams) => {
        const newSearchFilter =
          filters.length > 0 ? encodeFiltersForUrl(filters) : undefined;

        if (newSearchFilter !== search.filters) {
          navigate({
            search: (prev) => ({
              ...prev,
              filters: newSearchFilter,
            }),
            replace: true,
          });
        }
      }, 300),
    [navigate],
  );

  useEffect(() => {
    debouncedUrlUpdate(columnFilters, searchParams);
  }, [columnFilters, searchParams, debouncedUrlUpdate]);

  const columns: ColumnDef<EventUser>[] = useMemo(
    () => [
      {
        id: "avatar",
        header: "Avatar",
        cell: ({ row }) => {
          const avatarUrl = row.original.image;
          return avatarUrl ? (
            <img
              src={avatarUrl}
              alt="Avatar"
              className="h-10 w-10 rounded-full object-cover"
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
      },
      {
        accessorKey: "email",
        header: "Email",
        filterFn: fuzzyTextFilterFn,
      },
      {
        accessorKey: "event_role",
        header: "Role",
        cell: ({ row }) => {
          const role = row.original.event_role;
          if (!role) return "N/A";

          return <RoleBadge role={role} />;
        },
        filterFn: fuzzyTextFilterFn,
      },
      {
        id: "actions",
        header: "Actions",
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
              <UserSideDrawer user={row.original} />
            </DialogTrigger>
          );
        },
      },
    ],
    [],
  );

  const table = useReactTable({
    columns,
    data: data ?? fallbackData,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    state: { columnFilters },
    onColumnFiltersChange: setColumnFilters,
  });

  return (
    <div>
      <div>{eventId}</div>
      <table className="w-full">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="text-left px-4 py-2">
                  <div className="flex flex-col">
                    <div className="font-medium">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                    </div>
                    {header.column.getCanFilter() && (
                      <TextField
                        type="text"
                        value={(header.column.getFilterValue() ?? "") as string}
                        aria-label="Search"
                        onChange={(e) => header.column.setFilterValue(e)}
                        placeholder="Search..."
                        className="mt-2"
                      />
                    )}
                  </div>
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row, i) => (
            <tr key={row.id} className={i % 2 === 0 ? "bg-surface" : undefined}>
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id} className="p-4">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default UserTable;
