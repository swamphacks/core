import { useFieldContext } from "@/components/Form/formContext";
import { RadioGroup, type RadioGroupProps } from "@/components/ui/Radio";

export default function RadioField(props: RadioGroupProps) {
  const field = useFieldContext();

  return <RadioGroup onChange={(val) => field.handleChange(val)} {...props} />;
}
