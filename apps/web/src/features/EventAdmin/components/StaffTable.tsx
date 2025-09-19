import { Button } from "@/components/ui/Button";
import TablerTrash from "~icons/tabler/trash";

import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";
import RoleBadge from "./RoleBadge";
import type { StaffUser } from "@/features/PlatformAdmin/EventManager/hooks/useEventStaffUsers";

const columns: ColumnDef<StaffUser>[] = [
  {
    id: "avatar",
    header: "Avatar",
    cell: ({ row }) => {
      const avatarUrl = row.original.image;

      //TODO: Replace with actual avatar component
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
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "event_role",
    header: "Role",
    cell: ({ row }) => {
      const role = row.original.event_role;
      return <RoleBadge role={role} />;
    },
  },
  {
    id: "actions",
    header: "Actions",
    cell: ({ row }) => {
      const eventRole = row.original.event_role;

      return (
        <Button
          isDisabled={eventRole === "admin"}
          variant="danger"
          className="aspect-square p-2"
        >
          <TablerTrash className="h-4 w-4" />
        </Button>
      );
    },
  },
];

const fallbackData: StaffUser[] = [];

interface Props {
  data?: StaffUser[];
}

const StaffTable = ({ data }: Props) => {
  const table = useReactTable({
    columns,
    data: data ?? fallbackData,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div>
      <table className="w-full">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id} className="text-left px-4 py-2">
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
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

export default StaffTable;
