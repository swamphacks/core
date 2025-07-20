import { useFieldContext } from "@/components/Form/formContext";
import {
  MultiSelect,
  type MultiSelectProps,
} from "@/components/ui/MultiSelect";

export default function MultiSelectField(props: MultiSelectProps) {
  const field = useFieldContext();

  return (
    <MultiSelect
      onChange={(data) => field.setValue(data.map((item) => item.value))}
      errors={field.getMeta().errors.map((error) => error.message)}
      {...props}
    />
  );
}
