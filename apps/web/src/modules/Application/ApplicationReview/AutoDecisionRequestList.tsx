import { Button } from "@/components/ui/Button";
import { Popover } from "@/components/ui/Popover";
import { Sheet } from "@/components/ui/Sheet";
import { Table } from "@/components/ui/Table";
import ApplicationViewer from "@/modules/Application/ApplicationViewer";
import { useUpdateAutoDecisionRequest } from "@/modules/Application/hooks/useAutoDecisionRequests";
import useParsedForm from "@/modules/Application/hooks/useParsedForm";
import {
  useSearchAutoDecisionRequests,
  type SearchAutoDecisionRequestsResponse,
} from "@/modules/Application/hooks/useSearchAutoDecisionRequests";
import {
  type ColumnDef,
  useReactTable,
  getCoreRowModel,
  getFilteredRowModel,
} from "@tanstack/react-table";
import { useState, useMemo } from "react";
import { DialogTrigger } from "react-aria-components";

interface AutoDecisionRequestListProps {
  searchInput: string;
  approvalFilter: string;
  decisionFilter: string;
}

export default function AutoDecisionRequestList({
  searchInput,
  approvalFilter,
  decisionFilter,
}: AutoDecisionRequestListProps) {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 10,
  });

  const searchAutoDecisionRequestsData = useSearchAutoDecisionRequests(
    pagination.pageSize,
    pagination.pageIndex,
    searchInput,
    approvalFilter,
    decisionFilter,
  );

  const updateRequest = useUpdateAutoDecisionRequest();
  const parsedForm = useParsedForm();

  const requestRows =
    searchAutoDecisionRequestsData.data?.autoDecisionRequests ?? [];

  type AutoDecisionRequest =
    SearchAutoDecisionRequestsResponse["autoDecisionRequests"][number];

  const columns: ColumnDef<AutoDecisionRequest>[] = useMemo(
    () => [
      {
        id: "applicant",
        header: "Applicant",
        accessorKey: "user.name",
        size: 200,
        cell: ({ row }) => {
          const avatarUrl = row.original.user.image;
          return (
            <div className="flex items-center gap-2">
              {avatarUrl ? (
                <img
                  src={avatarUrl}
                  alt={"user avatar"}
                  className="size-8 rounded-full object-cover"
                />
              ) : (
                <div className="size-8 rounded-full bg-gray-300 dark:bg-neutral-700 flex items-center justify-center">
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
        id: "reviewer",
        header: "Reviewer",
        accessorKey: "reviewer.name",
        size: 200,
        cell: ({ row }) => {
          const avatarUrl = row.original.reviewer.image;
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
                {row.original.reviewer.name}
              </span>
            </div>
          );
        },
      },
      {
        header: "Application",
        size: 90,
        cell: ({ row }) => (
          <DialogTrigger>
            <Button variant="secondary" size="sm" className="h-8">
              Open
            </Button>
            <Sheet sheetClassName="w-full sm:w-160 lg:w-200">
              <ApplicationViewer
                parsedForm={parsedForm!}
                applicationId={row.original.applicationId}
              />
            </Sheet>
          </DialogTrigger>
        ),
      },
      {
        header: "Decision",
        accessorKey: "requestedDecision",
        enableGlobalFilter: false,
        size: 130,
        cell: ({ row }) =>
          row.original.requestedDecision === "auto_accept"
            ? "Auto Accept"
            : "Auto Reject",
      },
      {
        header: "Justification",
        enableGlobalFilter: false,
        size: 110,
        cell: ({ row }) => (
          <DialogTrigger>
            <Button variant="secondary" size="sm" className="h-8">
              View
            </Button>
            <Popover>
              <div className="p-2">
                <p className="max-w-60 text-wrap">
                  {row.original.justification}
                </p>
              </div>
            </Popover>
          </DialogTrigger>
        ),
      },
      {
        header: "Created At",
        size: 200,
        cell: ({ row }) => new Date(row.original.createdAt).toLocaleString(),
      },
      {
        header: "Actions",
        id: "actions",
        size: 150,
        accessorFn: (row) =>
          row.decidedBy ? (row.approved ? "approved" : "denied") : "pending",
        enableGlobalFilter: false,
        cell: ({ row }) => {
          const isResolved = row.original.decidedBy !== null;

          return isResolved ? (
            <span className="text-sm text-text-secondary">
              {row.original.approved ? "Approved" : "Denied"}
            </span>
          ) : (
            <div className="flex gap-2">
              <Button
                className="h-8 px-2"
                size="sm"
                onClick={() =>
                  updateRequest.mutateAsync({
                    requestId: row.id,
                    approved: true,
                  })
                }
                isDisabled={updateRequest.isPending}
              >
                Approve
              </Button>
              <Button
                className="h-8 px-2"
                size="sm"
                variant="danger"
                onClick={() =>
                  updateRequest.mutateAsync({
                    requestId: row.id,
                    approved: false,
                  })
                }
                isDisabled={updateRequest.isPending}
              >
                Deny
              </Button>
            </div>
          );
        },
      },
      {
        header: "Decided By",
        accessorKey: "decidedBy",
        enableGlobalFilter: false,
        size: 180,
        cell: ({ row }) => (
          <span className="text-sm inline-block max-w-40 font-medium truncate">
            {row.original.decidedBy?.name}
          </span>
        ),
      },
    ],
    [parsedForm],
  );

  const table = useReactTable({
    globalFilterFn: "includesString",
    columns,
    data: requestRows,
    rowCount: searchAutoDecisionRequestsData.data?.count,
    state: { pagination, globalFilter: searchInput },
    getCoreRowModel: getCoreRowModel(),
    onPaginationChange: setPagination,
    getFilteredRowModel: getFilteredRowModel(),
    manualPagination: true,
    manualFiltering: true,
  });

  if (searchAutoDecisionRequestsData.isLoading) {
    return <p>Loading requests....</p>;
  }

  if (!searchAutoDecisionRequestsData.data) {
    return <p>Unable to load requests.</p>;
  }

  return requestRows.length === 0 ? (
    <div className="mt-2 min-w-[1000px]">
      <p className="text-text-secondary">No requests found.</p>
    </div>
  ) : (
    <Table
      className="max-h-[78vh] overflow-auto w-full"
      headerClassName="text-text-secondary text-sm"
      table={table}
    />
  );
}
