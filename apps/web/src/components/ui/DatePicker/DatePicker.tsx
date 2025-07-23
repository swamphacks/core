import {
  DatePicker as AriaDatePicker,
  type DatePickerProps as AriaDatePickerProps,
  type DateValue,
  type ValidationResult,
} from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { Calendar } from "@/components/ui/Calendar";
import { DateInput } from "@/components/ui/DateField";
import { Dialog } from "@/components/ui/Dialog";
import {
  Description,
  FieldError,
  FieldGroup,
  Label,
} from "@/components/ui/Field";
import { Popover } from "@/components/ui/Popover";
import { composeTailwindRenderProps } from "@/components/ui/utils";
import TablerCalendar from "~icons/tabler/calendar";

export interface DatePickerProps<T extends DateValue>
  extends AriaDatePickerProps<T> {
  label?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
}

export function DatePicker<T extends DateValue>({
  label,
  description,
  errorMessage,
  ...props
}: DatePickerProps<T>) {
  return (
    <AriaDatePicker
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group flex flex-col gap-1",
      )}
    >
      {label && <Label isRequired={props.isRequired}>{label}</Label>}
      <FieldGroup className="min-w-[208px] w-auto border-1 rounded-sm items-center bg-surface">
        <DateInput className="flex-1 min-w-[150px] px-2 py-1.5 text-sm" />
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
      {description && <Description>{description}</Description>}
      <FieldError>{errorMessage}</FieldError>
      <Popover>
        <Dialog>
          <Calendar />
        </Dialog>
      </Popover>
    </AriaDatePicker>
  );
}
