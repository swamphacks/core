import { useFieldContext } from "@/components/Form/formContext";
import { Select, type SelectProps } from "@/components/ui/Select";

export default function SelectField<T extends object>(props: SelectProps<T>) {
  const field = useFieldContext();

  return (
    <Select onSelectionChange={(val) => field.handleChange(val)} {...props} />
  );
}
