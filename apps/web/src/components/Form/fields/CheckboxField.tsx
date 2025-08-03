import { useFieldContext } from "@/components/Form/formContext";
import {
  CheckboxGroup,
  type CheckboxGroupProps,
} from "@/components/ui/Checkbox";

export default function CheckboxField(props: CheckboxGroupProps) {
  const field = useFieldContext();

  return (
    <CheckboxGroup
      onChange={(val) => field.handleChange(val)}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
