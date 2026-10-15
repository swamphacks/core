import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Table } from "@/components/ui/Table";
import {
  type ColumnDef,
  useReactTable,
  getCoreRowModel,
} from "@tanstack/react-table";
import { useMemo } from "react";
import { DialogTrigger } from "react-aria-components";
import TablerPlus from "~icons/tabler/plus";
import { DeleteRedeemableModal } from "./DeleteRedeemableModal";
import { RedeemableFormModal } from "./RedeemableFormModal";
import {
  redeemableTypeInfo,
  useRedeemables,
  type Redeemable,
} from "./hooks/useRedeemables";

export default function RedeemablesPage() {
  const redeemables = useRedeemables();
  // Only list items the NFC app can scan, so staff don't see swag or drink
  // items that no app tab shows.
  const rows = useMemo(
    () => (redeemables.data ?? []).filter((r) => redeemableTypeInfo(r.type)),
    [redeemables.data],
  );

  const columns: ColumnDef<Redeemable>[] = useMemo(
    () => [
      {
        id: "name",
        header: "Name",
        accessorKey: "name",
        minSize: 220,
        cell: ({ row }) => (
          <span className="text-sm font-medium">{row.original.name}</span>
        ),
      },
      {
        id: "type",
        header: "Type",
        accessorKey: "type",
        size: 260,
        cell: ({ row }) => {
          const info = redeemableTypeInfo(row.original.type);
          return (
            <div className="flex items-center gap-2">
              <Badge>{info?.name ?? row.original.type}</Badge>
              <span className="text-xs text-text-secondary">
                {info?.appTab} tab
              </span>
            </div>
          );
        },
      },
      {
        id: "max_user_amount",
        header: "Limit per hacker",
        accessorKey: "max_user_amount",
        size: 140,
      },
      {
        id: "total_redeemed",
        header: "Times redeemed",
        accessorKey: "total_redeemed",
        size: 140,
      },
      {
        id: "actions",
        header: "",
        size: 180,
        cell: ({ row }) => (
          <div className="flex gap-2">
            <DialogTrigger>
              <Button variant="secondary" size="sm" className="h-8">
                Edit
              </Button>
              <RedeemableFormModal redeemable={row.original} />
            </DialogTrigger>
            <DialogTrigger>
              <Button variant="danger" size="sm" className="h-8">
                Delete
              </Button>
              <DeleteRedeemableModal redeemable={row.original} />
            </DialogTrigger>
          </div>
        ),
      },
    ],
    [],
  );

  const table = useReactTable({
    columns,
    data: rows,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div>
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold">Redeemables</h1>
          <p className="text-sm text-text-secondary">
            Meals and T-shirts staff hand out by scanning NFC tags. Each scan
            counts toward the hacker's limit for that item.
          </p>
        </div>
        <DialogTrigger>
          <Button variant="primary" className="w-fit">
            <TablerPlus className="w-4 h-4" />
            New Item
          </Button>
          <RedeemableFormModal />
        </DialogTrigger>
      </div>

      <div className="mt-5 p-2 rounded-lg bg-surface border border-border flex flex-col">
        {redeemables.isLoading ? (
          <p className="p-2">Loading items...</p>
        ) : redeemables.isError ? (
          <p className="p-2">Unable to load items.</p>
        ) : rows.length === 0 ? (
          <p className="p-2 text-text-secondary">
            No items yet. Create one for each meal or T-shirt giveaway.
          </p>
        ) : (
          <Table
            className="max-h-[78vh] overflow-auto w-full"
            headerClassName="text-text-secondary text-sm"
            table={table}
            showPagination={false}
          />
        )}
      </div>
    </div>
  );
}
