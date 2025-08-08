// components/event-staff/StaffTable.tsx
import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";
import type { User } from "@/lib/openapi/types";

export function StaffTable({
  data,
  columns,
  isLoading,
  isError,
}: {
  data: User[];
  columns: ColumnDef<User>[];
  isLoading: boolean;
  isError: boolean;
}) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const colSpan = table.getHeaderGroups()[0]?.headers.length || 1;

  return (
    <div className="overflow-x-auto max-h-48 overflow-y-auto border rounded-md border-input-border">
      <table className="w-full border-collapse">
        <thead className="bg-surface">
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th
                  key={header.id}
                  className="border border-input-border px-4 py-2 text-left text-sm font-semibold"
                >
                  {flexRender(
                    header.column.columnDef.header,
                    header.getContext(),
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {isLoading ? (
            <tr>
              <td
                colSpan={colSpan}
                className="text-center py-4 text-sm text-muted-foreground"
              >
                Loading...
              </td>
            </tr>
          ) : isError ? (
            <tr>
              <td
                colSpan={colSpan}
                className="text-center py-4 text-sm text-red-600"
              >
                Error loading data.
              </td>
            </tr>
          ) : table.getRowModel().rows.length === 0 ? (
            <tr>
              <td
                colSpan={colSpan}
                className="text-center py-4 text-sm text-muted-foreground"
              >
                No data available.
              </td>
            </tr>
          ) : (
            table.getRowModel().rows.map((row) => (
              <tr key={row.id} className="even:bg-surface">
                {row.getVisibleCells().map((cell) => (
                  <td
                    key={cell.id}
                    className="border px-4 py-2 text-sm border-input-border"
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
