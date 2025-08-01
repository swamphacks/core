import { useFieldContext } from "@/components/Form/formContext";
import { DatePicker, type DatePickerProps } from "@/components/ui/DatePicker";
import { getLocalTimeZone } from "@internationalized/date";
import type { DateValue } from "react-aria-components";

export default function DatePickerField<T extends DateValue>(
  props: DatePickerProps<T>,
) {
  const field = useFieldContext();

  return (
    <DatePicker
      onChange={(val) => field.handleChange(val?.toDate(getLocalTimeZone()))}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
