import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
  type ColumnDef,
  type OnChangeFn,
} from "@tanstack/react-table";
import {
  MultiSelect,
  type MultiSelectProps,
} from "@/components/ui/MultiSelect";
import { TextField } from "@/components/ui/TextField";
import { Button } from "../Button";
import TablerChevronUp from "~icons/tabler/chevron-up";
import TablerChevronDown from "~icons/tabler/chevron-down";
import TablerArrowsSort from "~icons/tabler/arrows-sort";

export interface ColumnFilter {
  id: string;
  value: unknown;
}
export type ColumnFiltersState = ColumnFilter[];

export type ColumnSort = {
  id: string;
  desc: boolean;
};

export type SortingState = ColumnSort[];

interface TableProps<TData> {
  data: TData[];
  columns: ColumnDef<TData>[];
  columnFilters: ColumnFiltersState;
  onColumnFiltersChange: OnChangeFn<ColumnFiltersState>;
  sorting: SortingState;
  onSortingChange: OnChangeFn<SortingState>;
  fallbackData?: TData[];
}

export function Table<TData>({
  data,
  columns,
  columnFilters,
  onColumnFiltersChange,
  sorting,
  onSortingChange,
  fallbackData = [],
}: TableProps<TData>) {
  const table = useReactTable({
    columns,
    data: data ?? fallbackData,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    state: { columnFilters, sorting },
    onColumnFiltersChange: onColumnFiltersChange,
    onSortingChange: onSortingChange,
  });

  return (
    <table className="w-full table-fixed">
      <thead>
        {table.getHeaderGroups().map((headerGroup) => (
          <tr key={headerGroup.id}>
            {headerGroup.headers.map((header) => {
              const sortState = header.column.getIsSorted();
              return (
                <th
                  key={header.id}
                  className={`text-left px-4 py-2`}
                  style={{
                    width: header.getSize(),
                    minWidth: header.column.columnDef.minSize,
                  }}
                >
                  <div className="flex flex-col">
                    <div className="font-medium">
                      {flexRender(
                        header.column.columnDef.header,
                        header.getContext(),
                      )}
                    </div>
                    <div className="flex w-full items-center justify-between gap-4 mb-2">
                      <div className="flex-1">
                        {header.column.getCanFilter() &&
                          header.column.columnDef.meta?.filterType ==
                            "text" && (
                            <TextField
                              type="text"
                              value={
                                (header.column.getFilterValue() ?? "") as string
                              }
                              aria-label="Search"
                              onChange={(e) => header.column.setFilterValue(e)}
                              placeholder="Search..."
                              className="mt-2"
                            />
                          )}
                        {header.column.getCanFilter() &&
                          header.column.columnDef.meta?.filterType ===
                            "select" &&
                          (() => {
                            const allOptions =
                              (header.column.columnDef.meta
                                ?.filterOptions as MultiSelectProps["options"]) ??
                              [];
                            const filterValue =
                              (header.column.getFilterValue() ??
                                []) as string[];
                            const selectedOptions = allOptions.filter((opt) =>
                              filterValue.includes(opt.value),
                            );

                            return (
                              <div className="mt-3">
                                <MultiSelect
                                  name={header.column.id}
                                  label=""
                                  options={allOptions}
                                  value={selectedOptions}
                                  onChange={(selected) => {
                                    const newValues = selected.map(
                                      (opt) => opt.value,
                                    );
                                    header.column.setFilterValue(
                                      newValues.length > 0
                                        ? newValues
                                        : undefined,
                                    );
                                  }}
                                />
                              </div>
                            );
                          })()}
                      </div>
                      {header.column.getCanSort() && (
                        <div className="mt-3">
                          <Button
                            variant="secondary" // Use a subtle button style
                            onClick={header.column.getToggleSortingHandler()}
                            className="p-3 h-10" // Make it small
                            aria-label={
                              sortState === "asc"
                                ? "Sorted ascending. Press to sort descending."
                                : sortState === "desc"
                                  ? "Sorted descending. Press to clear sort."
                                  : "Not sorted. Press to sort ascending."
                            }
                          >
                            {/* Conditionally render the correct icon */}
                            {sortState === "asc" ? (
                              <TablerChevronUp className="h-4 w-4" />
                            ) : sortState === "desc" ? (
                              <TablerChevronDown className="h-4 w-4" />
                            ) : (
                              <TablerArrowsSort className="h-4 w-4 opacity-30" />
                            )}
                          </Button>
                        </div>
                      )}
                    </div>
                  </div>
                </th>
              );
            })}
          </tr>
        ))}
      </thead>
      <tbody>
        {table.getRowModel().rows.map((row, i) => (
          <tr key={row.id} className={i % 2 === 0 ? "bg-surface" : undefined}>
            {row.getVisibleCells().map((cell) => {
              return (
                <td
                  key={cell.id}
                  className={`
                    p-4
                  `}
                >
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              );
            })}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
