import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  getPaginationRowModel,
  useReactTable,
  type ColumnDef,
  type OnChangeFn,
} from "@tanstack/react-table";
import {
  MultiSelect,
  type MultiSelectProps,
} from "@/components/ui/MultiSelect";
import { Select } from "@/components/ui/Select";
import { TextField } from "@/components/ui/TextField";
import { Button } from "../Button";
import TablerChevronUp from "~icons/tabler/chevron-up";
import TablerChevronDown from "~icons/tabler/chevron-down";
import TablerChevronRight from "~icons/tabler/chevron-right";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import TablerChevronsRight from "~icons/tabler/chevrons-right";
import TablerChevronsLeft from "~icons/tabler/chevrons-left";
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

export type PaginationState = {
  pageIndex: number;
  pageSize: number;
};

interface TableProps<TData> {
  data: TData[];
  columns: ColumnDef<TData>[];
  columnFilters: ColumnFiltersState;
  onColumnFiltersChange: OnChangeFn<ColumnFiltersState>;
  sorting: SortingState;
  onSortingChange: OnChangeFn<SortingState>;
  pagination: PaginationState;
  onPaginationChange: OnChangeFn<PaginationState>;
  fallbackData?: TData[];
}

export function Table<TData>({
  data,
  columns,
  columnFilters,
  onColumnFiltersChange,
  sorting,
  onSortingChange,
  pagination,
  onPaginationChange,
  fallbackData = [],
}: TableProps<TData>) {
  const table = useReactTable({
    columns,
    data: data ?? fallbackData,
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    state: { columnFilters, sorting, pagination },
    onColumnFiltersChange: onColumnFiltersChange,
    onSortingChange: onSortingChange,
    onPaginationChange: onPaginationChange,
  });

  return (
    <div>
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
                                  (header.column.getFilterValue() ??
                                    "") as string
                                }
                                aria-label="Search"
                                onChange={(e) =>
                                  header.column.setFilterValue(e)
                                }
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
      <div className="h-2" />
      <div className="flex items-center gap-2">
        <Button
          className="rounded p-1"
          variant="secondary"
          onClick={() => table.firstPage()}
          isDisabled={!table.getCanPreviousPage()}
        >
          <TablerChevronsLeft></TablerChevronsLeft>
        </Button>
        <Button
          className="rounded p-1"
          variant="secondary"
          onClick={() => table.previousPage()}
          isDisabled={!table.getCanPreviousPage()}
        >
          <TablerChevronLeft></TablerChevronLeft>
        </Button>
        <Button
          className="rounded p-1"
          variant="secondary"
          onClick={() => table.nextPage()}
          isDisabled={!table.getCanNextPage()}
        >
          <TablerChevronRight></TablerChevronRight>
        </Button>
        <Button
          className="rounded p-1"
          variant="secondary"
          onClick={() => table.lastPage()}
          isDisabled={!table.getCanNextPage()}
        >
          <TablerChevronsRight></TablerChevronsRight>
        </Button>
        <span className="flex items-center gap-1">
          <div>Page</div>
          <strong>
            {table.getState().pagination.pageIndex + 1} of{" "}
            {table.getPageCount().toLocaleString()}
          </strong>
        </span>
        {/* <span className="flex items-center gap-1">
          | Go to page:
          <input
            type="number"
            min="1"
            max={table.getPageCount()}
            defaultValue={table.getState().pagination.pageIndex + 1}
            onChange={(e) => {
              const page = e.target.value ? Number(e.target.value) - 1 : 0;
              table.setPageIndex(page);
            }}
            className="border p-1 rounded w-16"
          />
        </span> */}
        <Select
          aria-label="Select Page Size"
          selectedKey={String(table.getState().pagination.pageSize)}
          items={[5, 10, 20, 50, 100, 1000].map((size) => ({
            id: String(size),
            name: `Row Size: ${size}`,
          }))}
          onSelectionChange={(key) => {
            if (key) {
              table.setPageSize(Number(key));
            }
          }}
          children={null}
        ></Select>
      </div>
    </div>
  );
}
