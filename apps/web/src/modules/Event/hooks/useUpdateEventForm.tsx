// useEventForm.ts
import { useForm } from "@tanstack/react-form";
import { useUpdateEvent } from "../hooks/useUpdateEvent";
import { toast } from "react-toastify";
import { pickDirty } from "@/utils/formHelper";
import { EventSchema, type Event } from "../schemas/event";
import { toCalendarDateTime } from "@/utils/date";
import { z } from "zod";
import { CalendarDateTimeSchema } from "@/utils/customSchemas";
import { omit } from "@/utils/object";

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

export type UpdateEventFormType = z.infer<typeof UpdateEventSchema>;

export function useUpdateEventForm(event: Event) {
  const userTimeZone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  const {
    id,
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

  const { mutateAsync } = useUpdateEvent(id);

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
      if (!success) {
        console.error("Validation failed:", data);
        return;
      }

      const dirtyFields = pickDirty(value, form.state.fieldMeta);
      const formattedFields: Partial<Event> = {};

      // Convert event times
      if (dirtyFields.event_times) {
        formattedFields.start_time =
          dirtyFields.event_times.start.toDate(userTimeZone);
        formattedFields.end_time =
          dirtyFields.event_times.end.toDate(userTimeZone);
      }

      // Convert application times
      if (dirtyFields.application_times) {
        formattedFields.application_open =
          dirtyFields.application_times.start.toDate(userTimeZone);
        formattedFields.application_close =
          dirtyFields.application_times.end.toDate(userTimeZone);
      }

      // Convert optional fields
      if ("rsvp_deadline" in dirtyFields) {
        formattedFields.rsvp_deadline = dirtyFields.rsvp_deadline
          ? dirtyFields.rsvp_deadline.toDate(userTimeZone)
          : null;
      }
      if ("decision_release" in dirtyFields) {
        formattedFields.decision_release = dirtyFields.decision_release
          ? dirtyFields.decision_release.toDate(userTimeZone)
          : null;
      }

      // Merge other fields
      Object.assign(
        formattedFields,
        omit(
          dirtyFields,
          "event_times",
          "application_times",
          "rsvp_deadline",
          "decision_release",
        ),
      );

      // Submit
      await mutateAsync(formattedFields, {
        onSuccess: () =>
          toast.success("Event updated successfully", {
            position: "bottom-right",
          }),
        onError: () =>
          toast.error("Failed to update event", { position: "bottom-right" }),
      });

      form.reset(value);
    },
    validators: {
      onSubmit: UpdateEventSchema,
    },
  });

  return form;
}
