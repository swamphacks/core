import { CalendarDate, CalendarDateTime } from "@internationalized/date";

/**
 * Convert a string or Date to a CalendarDate.
 */
export function toCalendarDate(raw: string | Date): CalendarDate {
  const date = raw instanceof Date ? raw : new Date(raw);
  return new CalendarDate(
    date.getFullYear(),
    date.getMonth() + 1,
    date.getDate(),
  );
}

/**
 * Convert a string or Date to a CalendarDateTime with optional timezone.
 */
export function toCalendarDateTime(
  raw: string | Date,
  timezone: string = Intl.DateTimeFormat().resolvedOptions().timeZone,
): CalendarDateTime {
  const date = raw instanceof Date ? raw : new Date(raw);

  // Convert to the target timezone using Intl.DateTimeFormat
  const parts = new Intl.DateTimeFormat("en-US", {
    timeZone: timezone,
    year: "numeric",
    month: "numeric",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
    hour12: false,
  }).formatToParts(date);

  const getPart = (type: string) =>
    Number(parts.find((p) => p.type === type)?.value ?? 0);

  return new CalendarDateTime(
    getPart("year"),
    getPart("month"),
    getPart("day"),
    getPart("hour"),
    getPart("minute"),
    getPart("second"),
  );
}
