import type { CalendarDateTime, DateValue } from "@internationalized/date";
import z from "zod";

export const DateTimeSchema = z
  .custom<DateValue>()
  .transform((val) => val.toDate("UTC"));

export const CalendarDateTimeSchema = z
  .custom<CalendarDateTime>()
  .transform((val) => val.toDate("UTC"));
