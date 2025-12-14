import { Button } from "@/components/ui/Button";
import { TextField } from "@/components/ui/TextField";
import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import {
  Form,
  OverlayTriggerStateContext,
  Text,
  type DateRange,
} from "react-aria-components";
import z from "zod";
import { useCreateAdminEvent } from "../hooks/useCreateAdminEvent";
import { useFormErrors } from "@/components/Form";
import { Modal } from "@/components/ui/Modal";
import { DateRangePicker } from "@/components/ui/DateRangePicker";
import { useContext } from "react";
import { DateTimeSchema } from "@/utils/customSchemas";

const addEventSchema = z.object({
  eventName: z.string().min(1, "Event name is required"),
  eventDateRange: z.object(
    {
      start: DateTimeSchema,
      end: DateTimeSchema,
    },
    "Must specify a date range.",
  ),
  applicationDateRange: z.object(
    {
      start: DateTimeSchema,
      end: DateTimeSchema,
    },
    "Must specify a date range.",
  ),
});

export type AddEvent = z.infer<typeof addEventSchema>;

interface AddEventFormValues {
  eventName: string;
  eventDateRange: DateRange | null;
  applicationDateRange: DateRange | null;
}

function AddEventModal() {
  const state = useContext(OverlayTriggerStateContext)!;
  const { mutateAsync } = useCreateAdminEvent();

  const form = useForm<
    AddEventFormValues,
    undefined,
    undefined,
    undefined,
    undefined,
    undefined,
    FormValidateOrFn<AddEventFormValues>,
    undefined,
    undefined,
    undefined
  >({
    defaultValues: {
      eventName: "",
      eventDateRange: null,
      applicationDateRange: null,
    },
    onSubmit: async ({ value }) => {
      const { data, success } = addEventSchema.safeParse(value);
      if (!success) return; // This should never run

      await mutateAsync(data, {
        onSuccess: () => {
          state.close();
          form.reset();
        },
      });
    },
    validators: {
      onSubmit: addEventSchema,
    },
  });

  const errors = useFormErrors(form);

  return (
    <Modal title="Create an event" responsive="sheet" size="xl" padding="sm">
      <Form
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
        validationErrors={errors}
      >
        <div className="mt-4 flex flex-col gap-4">
          <form.Field name="eventName">
            {(field) => (
              <TextField
                label="Event Name"
                name={field.name}
                placeholder="SwampHacks X"
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

          <form.Field name="eventDateRange">
            {(field) => (
              <DateRangePicker
                startName={field.name}
                label="Select event dates"
                granularity="minute"
                className="w-full"
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

          <form.Field name="applicationDateRange">
            {(field) => (
              <DateRangePicker
                startName={field.name}
                label="Select application open & close"
                granularity="minute"
                className="w-full"
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

          <Text className="text-text-secondary mt-2">
            Organizers will be able to add more configurations in their portal.
          </Text>

          <Button
            isPending={form.state.isSubmitting}
            type="submit"
            variant="primary"
          >
            Create
          </Button>
        </div>
      </Form>
    </Modal>
  );
}

export { AddEventModal };
