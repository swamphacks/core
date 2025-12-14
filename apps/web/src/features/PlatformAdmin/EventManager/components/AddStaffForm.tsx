// components/event-staff/AddStaffForm.tsx
import { Form, useFormErrors } from "@/components/Form";
import { TextField } from "@/components/ui/TextField";
import { Select, SelectItem } from "@/components/ui/Select";
import { Button } from "@/components/ui/Button";
import { Group, Text } from "react-aria-components";
import { useForm } from "@tanstack/react-form";
import {
  assignStaffRoleSchema,
  type AssignStaffRole,
} from "../hooks/useAdminStaffActions";

export function AddStaffForm({
  onSubmit,
}: {
  onSubmit: (data: AssignStaffRole) => Promise<void>;
}) {
  const form = useForm({
    defaultValues: {
      email: "",
      role: "ADMIN",
    },
    onSubmit: async ({ value }) => {
      const { data, success } = assignStaffRoleSchema.safeParse(value);
      if (!success) return;
      await onSubmit(data);
      form.reset();
    },
    validators: {
      onSubmit: assignStaffRoleSchema,
    },
  });

  const errors = useFormErrors(form);

  const roleOptions = [
    { id: "admin", name: "Admin" },
    { id: "staff", name: "Staff" },
  ];

  return (
    <Form
      validationErrors={errors}
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <Group className="flex sm:flex-row flex-col gap-4 sm:gap-2 w-full">
        <div className="flex flex-row gap-2 w-full">
          <form.Field name="email">
            {(field) => (
              <TextField
                label="Add by email:"
                name={field.name}
                description="Use account email, not their preferred."
                placeholder="jhnsmith@gmail.com"
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="role">
            {(field) => (
              <Select
                label="Role"
                defaultSelectedKey="Staff"
                name={field.name}
                selectedKey={field.state.value}
                placeholder="Role"
                items={roleOptions}
                onSelectionChange={(value) => {
                  if (!value) return;
                  field.handleChange(value.valueOf() as string);
                }}
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                validationBehavior="aria"
              >
                {(item) => <SelectItem key={item.id}>{item.name}</SelectItem>}
              </Select>
            )}
          </form.Field>
        </div>

        <div className="flex flex-col justify-center">
          <Button
            className="h-fit px-6"
            variant="primary"
            type="submit"
            isPending={form.state.isSubmitting}
          >
            Add
          </Button>
        </div>
      </Group>

      {form.state.isSubmitSuccessful && (
        <Text className="text-green-400 mt-2">
          Successfully added staff member!
        </Text>
      )}
    </Form>
  );
}
