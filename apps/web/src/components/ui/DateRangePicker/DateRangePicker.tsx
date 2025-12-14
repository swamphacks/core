import {
  DateRangePicker as AriaDateRangePicker,
  RangeCalendar,
  type DateRangePickerProps as AriaDateRangePickerProps,
  type DateValue,
  type ValidationResult,
  Text,
  Heading,
  CalendarGrid,
  CalendarGridBody,
  CalendarGridHeader,
  CalendarCell,
  CalendarHeaderCell,
} from "react-aria-components";
import { FieldError, FieldGroup, Label } from "@/components/ui/Field";
import { composeTailwindRenderProps } from "@/components/ui/utils";
import TablerCalendar from "~icons/tabler/calendar";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import TablerChevronRight from "~icons/tabler/chevron-right";
import { Button } from "@/components/ui/Button";
import { useLocale } from "react-aria-components";
import { tv } from "tailwind-variants";
import { Dialog } from "@/components/ui/Dialog";
import { Popover } from "@/components/ui/Popover";
import { DateInput } from "../DateField";

const cellStyles = tv({
  base: "size-7 m-px text-sm cursor-default rounded-full flex items-center justify-center forced-color-adjust-none",
  variants: {
    isSelected: {
      false:
        "text-zinc-900 dark:text-zinc-200 hover:bg-gray-100 dark:hover:bg-zinc-700 pressed:bg-gray-200 dark:pressed:bg-zinc-600",
      true: "bg-blue-600 invalid:bg-red-600 text-white forced-colors:bg-[Highlight] forced-colors:invalid:bg-[Mark] forced-colors:text-[HighlightText]",
    },
    isDisabled: {
      true: "text-gray-300 dark:text-zinc-600 forced-colors:text-[GrayText]",
    },
  },
});

export interface DateRangePickerProps<T extends DateValue>
  extends AriaDateRangePickerProps<T> {
  label?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
}

export function DateRangePicker<T extends DateValue>({
  label,
  description,
  errorMessage,
  firstDayOfWeek,
  ...props
}: DateRangePickerProps<T>) {
  return (
    <AriaDateRangePicker
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group flex flex-col gap-1",
      )}
    >
      {label && <Label isRequired={props.isRequired}>{label}</Label>}
      <FieldGroup className="min-w-[208px] w-auto border-1 rounded-sm items-center bg-input-bg justify-evenly flex-nowrap gap-2">
        <DateInput
          slot="start"
          className="border-none min-w-[120px] flex-1 px-2 py-2 text-sm"
        />
        <span aria-hidden="true">-</span>
        <DateInput
          slot="end"
          className="border-none min-w-[150px] flex-1 px-2 py-2 text-sm"
        />

        <Button
          variant="icon"
          className="w-6 mr-1 rounded-sm outline-offset-0 px-0 py-1"
        >
          <TablerCalendar
            aria-hidden
            className="size-4 text-gray-600 dark:text-zinc-400 forced-colors:text-[ButtonText] group-disabled:text-gray-200 dark:group-disabled:text-zinc-600 forced-colors:group-disabled:text-[GrayText]"
          />
        </Button>
      </FieldGroup>
      {description && <Text slot="description">{description}</Text>}
      <FieldError>{errorMessage}</FieldError>

      <Popover>
        <Dialog>
          <RangeCalendar firstDayOfWeek={firstDayOfWeek}>
            <CalendarHeader />
            <CalendarGrid>
              <CalendarGridHeader>
                {(day) => (
                  <CalendarHeaderCell className="text-xs text-gray-500 font-semibold">
                    {day}
                  </CalendarHeaderCell>
                )}
              </CalendarGridHeader>
              <CalendarGridBody>
                {(date) => <CalendarCell className={cellStyles} date={date} />}
              </CalendarGridBody>
            </CalendarGrid>
          </RangeCalendar>
        </Dialog>
      </Popover>
    </AriaDateRangePicker>
  );
}

function CalendarHeader() {
  const { direction } = useLocale();

  return (
    <header className="flex items-center gap-1 pb-4 px-1 w-full">
      <Button variant="icon" slot="previous">
        {direction === "rtl" ? (
          <TablerChevronRight aria-hidden />
        ) : (
          <TablerChevronLeft aria-hidden />
        )}
      </Button>
      <Heading className="flex-1 font-semibold text-md text-center text-zinc-900 dark:text-zinc-200" />
      <Button variant="icon" slot="next">
        {direction === "rtl" ? (
          <TablerChevronLeft aria-hidden />
        ) : (
          <TablerChevronRight aria-hidden />
        )}
      </Button>
    </header>
  );
}
