import { Form, useFormErrors } from "@/components/Form";
import { useForm } from "@tanstack/react-form";
import z from "zod";
import { EventSchema, type Event } from "../schemas/event";
import { TextField } from "@/components/ui/TextField";
import { CalendarDateTimeSchema } from "@/utils/customSchemas";
import { DateRangePicker } from "@/components/ui/DateRangePicker";
import { toCalendarDateTime } from "@/utils/date";
import { DatePicker } from "@/components/ui/DatePicker";
import { Group } from "react-aria-components";

const UpdateEventSchema = EventSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
  start_time: true,
  end_time: true,
  application_open: true,
  application_close: true,
  rsvp_deadline: true,
  decision_release: true,
}).extend({
  event_times: z.object({
    start: CalendarDateTimeSchema,
    end: CalendarDateTimeSchema,
  }),
  application_times: z.object({
    start: CalendarDateTimeSchema,
    end: CalendarDateTimeSchema,
  }),
  rsvp_deadline: CalendarDateTimeSchema.nullable(),
  decision_release: CalendarDateTimeSchema.nullable(),
});

interface Props {
  event: Event;
}

const EventSettingsForm = ({ event }: Props) => {
  const {
    id, // eslint-disable-line @typescript-eslint/no-unused-vars
    created_at, // eslint-disable-line @typescript-eslint/no-unused-vars
    updated_at, // eslint-disable-line @typescript-eslint/no-unused-vars
    start_time,
    end_time,
    application_close,
    application_open,
    rsvp_deadline,
    decision_release,
    ...eventFields
  } = event;

  console.log("Event fields:", event);

  const form = useForm({
    defaultValues: {
      ...eventFields,
      event_times: {
        start: toCalendarDateTime(start_time),
        end: toCalendarDateTime(end_time),
      },
      application_times: {
        start: toCalendarDateTime(application_open),
        end: toCalendarDateTime(application_close),
      },
      rsvp_deadline: rsvp_deadline ? toCalendarDateTime(rsvp_deadline) : null,
      decision_release: decision_release
        ? toCalendarDateTime(decision_release)
        : null,
    },
    onSubmit: async ({ value }) => {
      const { data, success } = UpdateEventSchema.safeParse(value);
      if (!success) return;
      console.log("Form submitted:", data);
    },
    validators: {
      onSubmit: UpdateEventSchema,
    },
  });

  const errors = useFormErrors(form);

  return (
    <div className="max-w-lg">
      <Form
        validationErrors={errors}
        onSubmit={(e) => {
          e.preventDefault();
          e.stopPropagation();
          form.handleSubmit();
        }}
      >
        <div className="flex flex-col gap-4">
          <form.Field name="name">
            {(field) => (
              <TextField
                label="Event Title"
                name={field.name}
                isRequired
                placeholder="Ex. SwampHacks XI"
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

          <form.Field name="description">
            {(field) => (
              <TextField
                label="Description"
                name={field.name}
                placeholder="Ex. SwampHacks is Florida's largest student-run hackathon..."
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                textarea
                value={field.state.value ?? undefined}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="website_url">
            {(field) => (
              <TextField
                label="Event Website URL"
                name={field.name}
                placeholder="Ex. https://swamphacks.com"
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value ?? undefined}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="max_attendees">
            {(field) => (
              <TextField
                label="Maximum Attendees"
                name={field.name}
                placeholder="Enter the maximum number of attendees"
                className="flex-1"
                type="number"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={
                  field.state.value !== null && field.state.value !== undefined
                    ? String(field.state.value)
                    : ""
                }
                onChange={(value) => {
                  if (value === "" || value === null) {
                    field.handleChange(null); // allow null
                  } else {
                    const parsed = parseInt(value, 10);
                    field.handleChange(isNaN(parsed) ? null : parsed);
                  }
                }}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="location">
            {(field) => (
              <TextField
                label="Venue Name"
                name={field.name}
                placeholder="Ex. Marston Science Library"
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value ?? undefined}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="location_url">
            {(field) => (
              <TextField
                label="Venue Map URL"
                name={field.name}
                placeholder="Ex. https://maps.google.com/?q=Marston+Science+Library"
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value ?? undefined}
                onChange={(value) => field.handleChange(value)}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="event_times">
            {(field) => (
              <DateRangePicker
                startName={field.name}
                label="Event Times"
                granularity="minute"
                className="w-full"
                isRequired
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value}
                onChange={(value) => {
                  if (!value) return;
                  field.handleChange(value);
                }}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <form.Field name="application_times">
            {(field) => (
              <DateRangePicker
                startName={field.name}
                label="Application Times"
                granularity="minute"
                className="w-full"
                isRequired
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value}
                onChange={(value) => {
                  if (!value) return;
                  field.handleChange(value);
                }}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          <Group className="flex flex-row flex-wrap gap-4">
            <form.Field name="decision_release">
              {(field) => (
                <DatePicker
                  label="Decision Release Time"
                  name={field.name}
                  className="flex-1"
                  isDisabled={
                    form.state.isSubmitting &&
                    field.state.meta.errors.length === 0
                  }
                  value={field.state.value ?? undefined}
                  onChange={(value) => field.handleChange(value)}
                  validationBehavior="aria"
                />
              )}
            </form.Field>

            <form.Field name="rsvp_deadline">
              {(field) => (
                <DatePicker
                  label="Accepted RSVP Deadline"
                  name={field.name}
                  className="flex-1"
                  isDisabled={
                    form.state.isSubmitting &&
                    field.state.meta.errors.length === 0
                  }
                  value={field.state.value ?? undefined}
                  onChange={(value) => field.handleChange(value)}
                  validationBehavior="aria"
                />
              )}
            </form.Field>
          </Group>
        </div>
      </Form>
    </div>
  );
};

export default EventSettingsForm;
