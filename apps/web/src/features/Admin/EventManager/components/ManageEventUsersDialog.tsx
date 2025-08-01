import { Modal } from "@/components/ui/Modal";
import { Group, Heading } from "react-aria-components";
import {
  type ColumnDef,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { TextField } from "@/components/ui/TextField";
import { Select, SelectItem } from "@/components/ui/Select";
import { Button } from "@/components/ui/Button";
import { Form } from "@/components/Form";
import { useForm } from "@tanstack/react-form";

interface UserEventConfig {
  firstName?: string | null;
  lastName?: string | null;
  email: string;
  event_role: "STAFF" | "ADMIN";
}

const mockUsers: UserEventConfig[] = [
  {
    firstName: "Alex",
    lastName: "Wang",
    email: "alexanderwang@ufl.edu",
    event_role: "ADMIN",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  {
    firstName: "Phoenix",
    lastName: "Guttieriez",
    email: "phoenix@ufl.edu",
    event_role: "STAFF",
  },
  // ... rest omitted for brevity
];

const columns: ColumnDef<UserEventConfig>[] = [
  {
    id: "name",
    header: "Name",
    cell: ({ row }) => {
      const { firstName, lastName } = row.original;
      if (!firstName && !lastName) return "Unknown";
      return `${firstName ?? ""} ${lastName ?? ""}`.trim();
    },
  },
  {
    accessorKey: "email",
    header: "Email",
  },
  {
    accessorKey: "event_role",
    header: "Role",
  },
];

function ManageEventUsersDialog() {
  const table = useReactTable({
    data: mockUsers, // <-- Correct prop here
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  const form = useForm({
    defaultValues: {
      email: "",
      role: "admin",
    },
  });

  return (
    <Modal isDismissible size="xl">
      <div className="flex flex-col gap-4">
        <Heading className="text-lg" slot="title">
          Manage Users
        </Heading>

        {/* Table */}
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
              {table.getRowModel().rows.map((row) => (
                <tr key={row.id} className="even:bg-surface">
                  {row.getVisibleCells().map((cell) => (
                    <td
                      key={cell.id}
                      className="border px-4 py-2 text-sm border-input-border"
                    >
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </td>
                  ))}
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <Form>
          <Group className="flex flex-row gap-2 justify-center">
            <form.Field name="email">
              {(field) => (
                <TextField
                  label="Add by email:"
                  name={field.name}
                  description="Use account email, not their preferred."
                  placeholder="jhnsmith@gmail.com"
                  className="flex-1"
                />
              )}
            </form.Field>

            <Select selectedKey="admin" label="Role" placeholder="Role">
              <SelectItem key="admin">Admin</SelectItem>
              <SelectItem key="staff">Staff</SelectItem>
            </Select>

            <div className="flex flex-col justify-center">
              <Button className="h-fit" variant="primary">
                Add
              </Button>
            </div>
          </Group>
        </Form>
      </div>
    </Modal>
  );
}

export { ManageEventUsersDialog };
