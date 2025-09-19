import { Form, useFormErrors } from "@/components/Form";
import { type Event } from "../schemas/event";
import { TextField } from "@/components/ui/TextField";
import { DateRangePicker } from "@/components/ui/DateRangePicker";
import { DatePicker } from "@/components/ui/DatePicker";
import { Group } from "react-aria-components";
import { Switch } from "@/components/ui/Switch";
import { Button } from "@/components/ui/Button";
import { useUpdateEventForm } from "../hooks/useUpdateEventForm";

interface Props {
  event: Event;
}

const EventSettingsForm = ({ event }: Props) => {
  const form = useUpdateEventForm(event);

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
          {/* Event Title */}
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
                onChange={field.handleChange}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          {/* Description */}
          <form.Field name="description">
            {(field) => (
              <TextField
                label="Description"
                name={field.name}
                placeholder="Ex. SwampHacks is Florida's largest student-run hackathon..."
                className="flex-1"
                textarea
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={field.state.value ?? undefined}
                onChange={field.handleChange}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          {/* Website URL */}
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
                onChange={field.handleChange}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          {/* Maximum Attendees */}
          <form.Field name="max_attendees">
            {(field) => (
              <TextField
                label="Maximum Attendees"
                name={field.name}
                type="number"
                placeholder="Enter the maximum number of attendees"
                className="flex-1"
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                value={
                  field.state.value != null ? String(field.state.value) : ""
                }
                onChange={(value) => {
                  const parsed = value === "" ? null : parseInt(value, 10);
                  field.handleChange(isNaN(parsed as number) ? null : parsed);
                }}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          {/* Venue */}
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
                onChange={field.handleChange}
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
                onChange={field.handleChange}
                validationBehavior="aria"
              />
            )}
          </form.Field>

          {/* Event Times */}
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

          {/* Application Times */}
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

          {/* Decision Release & RSVP */}
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
                  onChange={field.handleChange}
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
                  onChange={field.handleChange}
                  validationBehavior="aria"
                />
              )}
            </form.Field>
          </Group>

          {/* Published Switch */}
          <form.Field name="is_published">
            {(field) => (
              <Switch
                size="xl"
                className="mt-4"
                name={field.name}
                isDisabled={
                  form.state.isSubmitting &&
                  field.state.meta.errors.length === 0
                }
                isSelected={field.state.value}
                onChange={field.handleChange}
              >
                {field.state.value ? "Event is published" : "Event is unpublished"}
              </Switch>
            )}
          </form.Field>

          {/* Save Button */}
          <form.Subscribe selector={(state) => state.isDirty}>
            {(isDirty) => (
              <>
                {isDirty && (
                  <p className="text-sm text-orange-500 mt-2">
                    You have unsaved changes
                  </p>
                )}
                <Button type="submit" isDisabled={!isDirty} className="mt-2">
                  Save
                </Button>
              </>
            )}
          </form.Subscribe>
        </div>
      </Form>
    </div>
  );
};

export default EventSettingsForm;
