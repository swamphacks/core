import { useFieldContext } from "@/components/Form/formContext";
import { ComboBox, type ComboBoxProps } from "@/components/ui/ComboBox";

export default function ComboBoxField<T extends { id: string; name: string }>(
  props: ComboBoxProps<T>,
) {
  const field = useFieldContext();

  return (
    <ComboBox
      onSelectionChange={(val) => field.handleChange(val)}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
