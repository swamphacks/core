import { useFieldContext } from "@/components/Form/formContext";
import { DatePicker, type DatePickerProps } from "@/components/ui/DatePicker";
import type { DateValue } from "react-aria-components";

export default function DatePickerField<T extends DateValue>(
  props: DatePickerProps<T>,
) {
  const field = useFieldContext();

  return (
    <DatePicker
      onChange={(val) => field.handleChange(val)}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
