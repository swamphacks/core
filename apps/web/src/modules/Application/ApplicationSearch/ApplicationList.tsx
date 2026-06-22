import { Button } from "@/components/ui/Button";
import { Sheet } from "@/components/ui/Sheet";
import { Table } from "@/components/ui/Table";
import ApplicationViewer from "@/modules/Application/ApplicationViewer";
import { useSearchApplications } from "@/modules/Application/hooks/useSearchApplications";
import useParsedForm from "@/modules/Application/hooks/useParsedForm";
import {
  type ColumnDef,
  useReactTable,
  getCoreRowModel,
} from "@tanstack/react-table";
import { useState, useMemo } from "react";
import { DialogTrigger } from "react-aria-components";

interface ApplicationListProps {
  searchInput: string;
}

export default function ApplicationList({ searchInput }: ApplicationListProps) {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  });

  const searchApplicationsData = useSearchApplications(
    pagination.pageSize,
    pagination.pageIndex,
    searchInput,
  );

  const parsedForm = useParsedForm();

  const applicationRows = searchApplicationsData.data?.applications ?? [];

  const columns: ColumnDef<(typeof applicationRows)[number]>[] = useMemo(
    () => [
      {
        id: "name",
        header: "Username",
        accessorKey: "name",
        minSize: 250,
        cell: ({ row }) => {
          const avatarUrl = row.original.user.image;
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
                {row.original.user.name}
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
              <ApplicationViewer
                parsedForm={parsedForm!}
                applicationId={row.original.id}
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
        cell: ({ row }) => (row.original.isEarly ? "Yes" : "No"),
      },
      {
        id: "submitted_at",
        header: "Submitted At",
        accessorKey: "submitted_at",
        enableGlobalFilter: false,
        size: 200,
        cell: ({ row }) =>
          row.original.submittedAt
            ? new Date(row.original.submittedAt).toLocaleString("en-US", {
                timeZone: "America/New_York",
              })
            : null,
      },
      {
        id: "email",
        header: "Email",
        accessorKey: "email",
        size: 300,
        cell: ({ row }) => row.original.user.email,
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
    rowCount: searchApplicationsData.data?.count,
    manualPagination: true,
    manualFiltering: true,
  });

  if (searchApplicationsData.isLoading) {
    return <p>Loading applications....</p>;
  }

  if (!searchApplicationsData.data) {
    return <p>Unable to load applications.</p>;
  }

  return applicationRows.length === 0 ? (
    <div className="mt-2 min-w-[1000px]">
      <p className="text-text-secondary">No applications found.</p>
    </div>
  ) : (
    <Table
      className="max-h-[78vh] overflow-auto w-full"
      headerClassName="text-text-secondary text-sm"
      table={applicationsTable}
    />
  );
}
