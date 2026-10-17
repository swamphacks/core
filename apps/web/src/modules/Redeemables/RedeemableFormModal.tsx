import { Button } from "@/components/ui/Button";
import { TextField } from "@/components/ui/TextField";
import { Select } from "@/components/ui/Select";
import { useForm, useStore, type FormValidateOrFn } from "@tanstack/react-form";
import { Form, OverlayTriggerStateContext } from "react-aria-components";
import z from "zod";
import { useFormErrors } from "@/components/Form";
import { Modal } from "@/components/ui/Modal";
import { useContext } from "react";
import {
  apiErrorMessage,
  REDEEMABLE_TYPES,
  redeemableTypeInfo,
  useCreateRedeemable,
  useUpdateRedeemable,
  type Redeemable,
  type RedeemableType,
} from "./hooks/useRedeemables";
import { showToast } from "@/lib/toast/toast";

const redeemableSchema = z.object({
  name: z.string().trim().min(1, "Name is required"),
  type: z.enum(["meal", "tshirt"], "Pick a type"),
  maxUserAmount: z
    .number("Enter a number")
    .int("Must be a whole number")
    .min(1, "Must be at least 1"),
});

interface RedeemableFormValues {
  name: string;
  type: RedeemableType | "";
  maxUserAmount: number;
}

interface RedeemableFormModalProps {
  // Edits this item when given; creates a new one otherwise.
  redeemable?: Redeemable;
}

export function RedeemableFormModal({ redeemable }: RedeemableFormModalProps) {
  const state = useContext(OverlayTriggerStateContext)!;
  const { mutateAsync: createRedeemable } = useCreateRedeemable();
  const { mutateAsync: updateRedeemable } = useUpdateRedeemable();
  const isEdit = redeemable !== undefined;

  const form = useForm<
    RedeemableFormValues,
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    FormValidateOrFn<RedeemableFormValues>,
    undefined,
    undefined,
    undefined
  >({
    defaultValues: {
      name: redeemable?.name ?? "",
      type: (redeemableTypeInfo(redeemable?.type ?? "")?.id ?? "") as
        | RedeemableType
        | "",
      maxUserAmount: redeemable?.max_user_amount ?? 1,
    },
    onSubmit: async ({ value }) => {
      const { data, success } = redeemableSchema.safeParse(value);
      if (!success) return;

      try {
        if (isEdit) {
          await updateRedeemable({ id: redeemable.id, body: data });
        } else {
          await createRedeemable(data);
        }
        showToast({
          title: "Success",
          message: isEdit
            ? `Saved ${data.name}.`
            : `Created ${data.name}. It will show up in the app after staff reopen it.`,
          type: "success",
        });
        state.close();
        form.reset();
      } catch (error) {
        showToast({
          title: "Error",
          message: await apiErrorMessage(
            error,
            isEdit
              ? "Failed to save the item. Please try again."
              : "Failed to create the item. Please try again.",
          ),
          type: "error",
        });
      }
    },
    validators: {
      onSubmit: redeemableSchema,
    },
  });

  const errors = useFormErrors(form);
  // Read through useStore so the modal re-renders when submitting ends;
  // form.state alone isn't reactive, which left the fields disabled after
  // a failed validation.
  const isSubmitting = useStore(form.store, (state) => state.isSubmitting);

  return (
    <Modal
      title={isEdit ? "Edit Item" : "New Item"}
      responsive="sheet"
      size="xl"
      padding="sm"
    >
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
                placeholder="e.g. Saturday Lunch"
                description="Use one item per meal or giveaway, e.g. separate items for Saturday and Sunday lunch."
                isDisabled={isSubmitting}
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
                isRequired
              />
            )}
          </form.Field>

          <form.Field name="type">
            {(field) => {
              const info = redeemableTypeInfo(field.state.value);
              return (
                <Select
                  label="Type"
                  name={field.name}
                  placeholder="Pick a type"
                  items={REDEEMABLE_TYPES}
                  value={field.state.value || null}
                  onChange={(value) => {
                    if (value != null) {
                      field.handleChange(value.toString() as RedeemableType);
                    }
                  }}
                  description={
                    info
                      ? `Shows in the app's ${info.appTab} tab.`
                      : "Sets which tab the item shows in on the NFC app."
                  }
                  isDisabled={isSubmitting}
                  validationBehavior="aria"
                  isRequired
                  children={null}
                />
              );
            }}
          </form.Field>

          <form.Field name="maxUserAmount">
            {(field) => (
              <TextField
                label="Limit per hacker"
                name={field.name}
                description="How many times one hacker can redeem this item."
                type="number"
                inputMode="numeric"
                isDisabled={isSubmitting}
                value={
                  Number.isNaN(field.state.value)
                    ? ""
                    : field.state.value.toString()
                }
                onChange={(value) => field.handleChange(parseInt(value, 10))}
                validationBehavior="aria"
                isRequired
              />
            )}
          </form.Field>

          <Button
            isPending={isSubmitting}
            type="submit"
            variant="primary"
            className="mt-4"
          >
            {isEdit ? "Save" : "Create Item"}
          </Button>
        </div>
      </Form>
    </Modal>
  );
}
