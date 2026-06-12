import { useMemo, useState } from "react";
import type { components } from "@/lib/openapi/schema";
import { Button } from "@/components/ui/Button";
import { useDebounce } from "@uidotdev/usehooks";
import {
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from "@tanstack/react-table";
import { Table } from "@/components/ui/Table";
import { DialogTrigger } from "react-aria-components";
import { Sheet } from "@/components/ui/Sheet";
import { Input } from "@/components/ui/Field";
import {
  useAllApplications,
  type AllApplicationsData,
} from "@/modules/Application/hooks/useAllApplications";
import useParsedForm from "@/modules/Application/hooks/useParsedForm";
import ApplicationResponsesViewer from "@/modules/Application/ApplicationResponsesViewer";

type ApplicationRow = components["schemas"]["ListAllApplicationsRow"];

export default function ApplicationSearchPage() {
  const [searchInput, setSearchInput] = useState("");
  const debouncedSearchInput = useDebounce(searchInput, 500);
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  });

  const allApplicationsData = useAllApplications(
    pagination.pageSize,
    pagination.pageIndex,
    debouncedSearchInput,
  );

  const parsedForm = useParsedForm();

  const applicationRows: AllApplicationsData["applications"] =
    allApplicationsData.data?.applications ?? [];

  const columns: ColumnDef<ApplicationRow>[] = useMemo(
    () => [
      {
        id: "name",
        header: "Username",
        accessorKey: "name",
        minSize: 250,
        cell: ({ row }) => {
          const avatarUrl = row.original.image;
          return (
            <div className="flex items-center gap-2">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={"user avatar"}
                  className="h-8 w-8 rounded-full object-cover"
                />
              ) : (
                <div className="h-10 w-10 rounded-full bg-gray-300 dark:bg-neutral-700 flex items-center justify-center">
                  <span className="text-gray-600 dark:text-neutral-400">
                    N/A
                  </span>
                </div>
              )}
              <span className="text-sm inline-block max-w-40 font-medium truncate">
                {row.original.name}
              </span>
            </div>
          );
        },
      },
      {
        id: "status",
        header: "Status",
        size: 130,
        accessorKey: "status",
        enableGlobalFilter: false,
        cell: ({ row }) => {
          switch (row.original.status) {
            case "under_review":
              return "Under Review";
            // Capitalize first letter
            default:
              return (
                row.original.status.charAt(0).toUpperCase() +
                row.original.status.slice(1)
              );
          }
        },
      },
      {
        id: "application",
        header: "Application",
        size: 110,
        cell: ({ row }) => (
          <DialogTrigger>
            <Button variant="secondary" size="sm" className="h-8">
              Open
            </Button>
            <Sheet sheetClassName="w-full sm:w-160 lg:w-200">
              <ApplicationResponsesViewer
                parsedForm={parsedForm!}
                applicationId={row.original.user_id}
              />
            </Sheet>
          </DialogTrigger>
        ),
      },
      {
        id: "is_early",
        header: "Is Early",
        accessorKey: "is_early",
        size: 100,
        cell: ({ row }) => (row.original.is_early ? "Yes" : "No"),
      },
      {
        id: "submitted_at",
        header: "Submitted At",
        accessorKey: "submitted_at",
        enableGlobalFilter: false,
        size: 200,
        cell: ({ row }) =>
          row.original.submitted_at
            ? new Date(row.original.submitted_at).toLocaleString("en-US", {
                timeZone: "America/New_York",
              })
            : null,
      },
      {
        id: "email",
        header: "Email",
        accessorKey: "email",
        size: 300,
        cell: ({ row }) => row.original.email,
      },
    ],
    [parsedForm],
  );

  const applicationsTable = useReactTable({
    globalFilterFn: "includesString",
    columns,
    data: applicationRows,
    state: { pagination, globalFilter: searchInput },
    getCoreRowModel: getCoreRowModel(),
    // getFilteredRowModel: getFilteredRowModel(),
    onPaginationChange: setPagination,
    // onGlobalFilterChange: setSearchInput,
    // autoResetPageIndex: true,
    rowCount: allApplicationsData.data?.count,
    manualPagination: true,
    manualFiltering: true,
  });

  if (allApplicationsData.isLoading) {
    return <p>Loading applications....</p>;
  }

  if (!allApplicationsData.data) {
    return <p>Unable to load applications.</p>;
  }

  return (
    <div>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Application search</h1>
          <p className="text-sm text-muted-foreground">
            Browse applications and inspect status for each registered user.
          </p>
        </div>
      </div>

      <div className="mt-5 p-2 rounded-md bg-surface max-w-fit flex flex-col">
        <div className="w-70">
          <Input
            aria-label="search input"
            className="rounded-md"
            placeholder="Search by username or email"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
          />
        </div>
        {applicationRows.length === 0 ? (
          <div className="mt-2 min-w-[1000px]">
            <p className="text-text-secondary">No applications found.</p>
          </div>
        ) : (
          <Table
            className="max-h-100 overflow-auto w-full"
            headerClassName="text-text-secondary text-sm"
            table={applicationsTable}
          />
        )}
      </div>
    </div>
  );
}
