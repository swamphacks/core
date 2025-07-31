import { Button } from "@/components/ui/Button";
import { DateRangePicker } from "@/components/ui/DateRangePicker/DateRangePicker";
import { Modal } from "@/components/ui/Modal/Modal";
import { TextField } from "@/components/ui/TextField";
import { useForm, type FormValidateOrFn } from "@tanstack/react-form";
import {
  Form,
  Text,
  type DateRange,
  type DateValue,
} from "react-aria-components";
import z from "zod";
import { useCreateAdminEvent } from "../hooks/useCreateAdminEvent";

const calenderDateTimeSchema = z
  .custom<DateValue>()
  .transform((val) => val.toDate("UTC"));

const addEventSchema = z.object({
  eventName: z.string().min(1, "Event name is required"),
  eventDateRange: z.object({
    start: calenderDateTimeSchema,
    end: calenderDateTimeSchema,
  }),
  applicationDateRange: z.object({
    start: calenderDateTimeSchema,
    end: calenderDateTimeSchema,
  }),
});

export type AddEvent = z.infer<typeof addEventSchema>;

interface AddEventFormValues {
  eventName: string;
  eventDateRange: DateRange | null;
  applicationDateRange: DateRange | null;
}

interface AddEventModalProps {
  onClose: () => void;
}

function AddEventModal({ onClose }: AddEventModalProps) {
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
          onClose();
          form.reset();
        },
      });
    },
    validators: {
      onSubmit: addEventSchema,
    },
  });
  return (
    <Modal title="Create an event" responsive="sheet" size="xl" padding="sm">
      <Form
        onSubmit={(e) => {
          e.preventDefault();
          form.handleSubmit();
        }}
      >
        <div className="mt-4 flex flex-col gap-4">
          <form.Field name="eventName">
            {(field) => (
              <TextField
                label="Event Name"
                placeholder="SwampHacks X"
                isRequired
                isDisabled={form.state.isSubmitting}
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
              />
            )}
          </form.Field>

          <form.Field name="eventDateRange">
            {(field) => (
              <DateRangePicker
                isRequired
                label="Select event dates"
                granularity="minute"
                className="w-full"
                isDisabled={form.state.isSubmitting}
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
                errorMessage={field.state.meta.errors[0]}
              />
            )}
          </form.Field>

          <form.Field name="applicationDateRange">
            {(field) => (
              <DateRangePicker
                isRequired
                label="Select application open & close"
                granularity="minute"
                className="w-full"
                isDisabled={form.state.isSubmitting}
                value={field.state.value}
                onChange={(value) => field.handleChange(value)}
                errorMessage={field.state.meta.errors[0]}
              />
            )}
          </form.Field>

          <Text className="text-text-secondary mt-2">
            Organizers will be able to add more configurations in their portal.
          </Text>

          <Button
            isPending={form.state.isSubmitting}
            onClick={() => console.log(form.state.errorMap)}
            type="submit"
            variant="primary"
            className="mt-0"
          >
            Create
          </Button>
        </div>
      </Form>
    </Modal>
  );
}

export { AddEventModal };
