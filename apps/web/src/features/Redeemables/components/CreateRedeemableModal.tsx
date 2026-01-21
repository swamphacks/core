import { Button } from "@/components/ui/Button";
import { TextField } from "@/components/ui/TextField";
import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import { Form, OverlayTriggerStateContext } from "react-aria-components";
import z from "zod";
import { useFormErrors } from "@/components/Form";
import { Modal } from "@/components/ui/Modal";
import { useContext } from "react";
import { useCreateRedeemable } from "../hooks/useRedeemables";
import { showToast } from "@/lib/toast/toast";

const createRedeemableSchema = z.object({
  name: z.string().min(1, "Name is required"),
  amount: z.number().int().min(1, "Amount must be at least 1"),
  maxUserAmount: z.number().int().min(1, "Max user amount must be at least 1"),
});

interface CreateRedeemableFormValues {
  name: string;
  amount: number;
  maxUserAmount: number;
}

interface CreateRedeemableModalProps {
  eventId: string;
}

export function CreateRedeemableModal({ eventId }: CreateRedeemableModalProps) {
  const state = useContext(OverlayTriggerStateContext)!;
  const { mutateAsync } = useCreateRedeemable(eventId);

  const form = useForm<
    CreateRedeemableFormValues,
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    FormValidateOrFn<CreateRedeemableFormValues>,
    undefined,
    undefined,
    undefined
  >({
    defaultValues: {
      name: "",
      amount: 1,
      maxUserAmount: 1,
    },
    onSubmit: async ({ value }) => {
      const { data, success } = createRedeemableSchema.safeParse(value);
      if (!success) return;

      try {
        await mutateAsync({
          name: data.name,
          amount: data.amount,
          max_user_amount: data.maxUserAmount,
        });
        showToast({
          title: "Success",
          message: "Redeemable created successfully!",
          type: "success",
        });
        state.close();
        form.reset();
      } catch (error) {
        showToast({
          title: "Error",
          message: "Failed to create redeemable. Please try again.",
          type: "error",
        });
        console.log("Error creating redeemable:", error);
      }
    },
    validators: {
      onSubmit: createRedeemableSchema,
    },
  });

  const errors = useFormErrors(form);

  return (
    <Modal title="Create Redeemable" responsive="sheet" size="xl" padding="sm">
      <Form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
        validationErrors={errors}
      >
        <div className="mt-4 flex flex-col gap-4">
          <form.Field name="name">
            {(field) => (
              <TextField
                label="Name"
                name={field.name}
                placeholder="Enter redeemable name"
                isDisabled={form.state.isSubmitting}
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
                isRequired
              />
            )}
          </form.Field>

          <form.Field name="amount">
            {(field) => (
              <TextField
                label="Total Stock"
                name={field.name}
                placeholder="Enter total stock amount"
                type="number"
                isDisabled={form.state.isSubmitting}
                value={field.state.value?.toString() || ""}
                onChange={(value) =>
                  field.handleChange(value ? parseInt(value) : 0)
                }
                validationBehavior="aria"
                isRequired
              />
            )}
          </form.Field>

          <form.Field name="maxUserAmount">
            {(field) => (
              <TextField
                label="Max Per User"
                name={field.name}
                placeholder="Enter maximum amount per user"
                type="number"
                isDisabled={form.state.isSubmitting}
                value={field.state.value?.toString() || ""}
                onChange={(value) =>
                  field.handleChange(value ? parseInt(value) : 0)
                }
                validationBehavior="aria"
                isRequired
              />
            )}
          </form.Field>

          <Button
            isPending={form.state.isSubmitting}
            type="submit"
            variant="primary"
            className="mt-4"
          >
            Create Redeemable
          </Button>
        </div>
      </Form>
    </Modal>
  );
}
