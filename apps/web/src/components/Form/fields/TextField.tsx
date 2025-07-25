import { useFieldContext } from "@/components/Form/formContext";
import {
  TextField as BaseTextField,
  type TextFieldProps,
} from "@/components/ui/TextField";

export default function TextField(props: TextFieldProps) {
  const field = useFieldContext();

  return (
    <BaseTextField
      onChange={(val) => field.handleChange(val)}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
