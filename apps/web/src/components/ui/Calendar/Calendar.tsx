import {
  Calendar as AriaCalendar,
  CalendarGridHeader as AriaCalendarGridHeader,
  type CalendarProps as AriaCalendarProps,
  CalendarCell,
  CalendarGrid,
  CalendarGridBody,
  CalendarHeaderCell,
  type DateValue,
  Heading,
  Text,
  useLocale,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import { Button } from "@/components/ui/Button";
import TablerChevronLeft from "~icons/tabler/chevron-left";
import TablerChevronRight from "~icons/tabler/chevron-right";

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

export interface CalendarProps<T extends DateValue>
  extends Omit<AriaCalendarProps<T>, "visibleDuration"> {
  errorMessage?: string;
}

export function Calendar<T extends DateValue>({
  errorMessage,
  ...props
}: CalendarProps<T>) {
  return (
    <AriaCalendar {...props}>
      <CalendarHeader />
      <CalendarGrid>
        <CalendarGridHeader />
        <CalendarGridBody>
          {(date) => <CalendarCell date={date} className={cellStyles} />}
        </CalendarGridBody>
      </CalendarGrid>
      {errorMessage && (
        <Text slot="errorMessage" className="text-sm text-red-600">
          {errorMessage}
        </Text>
      )}
    </AriaCalendar>
  );
}

export function CalendarHeader() {
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

export function CalendarGridHeader() {
  return (
    <AriaCalendarGridHeader>
      {(day) => (
        <CalendarHeaderCell className="text-xs text-gray-500 font-semibold">
          {day}
        </CalendarHeaderCell>
      )}
    </AriaCalendarGridHeader>
  );
}
