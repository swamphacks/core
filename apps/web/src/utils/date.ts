import { CalendarDate, CalendarDateTime } from "@internationalized/date";

export function toCalendarDate(date: Date): CalendarDate {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  return new CalendarDate(year, month, day);
}

export function toCalendarDateTime(date: Date): CalendarDateTime {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  const hour = date.getHours();
  const minute = date.getMinutes();
  const second = date.getSeconds();
  return new CalendarDateTime(year, month, day, hour, minute, second);
}
