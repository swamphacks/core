import { useFieldContext } from "@/components/Form/formContext";
import {
  TextField as BaseTextField,
  type TextFieldProps,
} from "@/components/ui/TextField";

export default function TextField(
  props: TextFieldProps & { showWordCount?: boolean },
) {
  const field = useFieldContext();

  if (props.textarea && props.showWordCount) {
    const countWords = (text: string) => {
      if (text === undefined || text === null) return 0;

      // Remove leading/trailing whitespace and split by one or more spaces
      const words = text.trim().split(/\s+/);
      // Filter out empty strings that might result from multiple spaces
      return words.filter((word) => word !== "").length;
    };

    return (
      <div>
        <BaseTextField
          onChange={(val) => field.handleChange(val)}
          onBlur={field.handleBlur}
          {...props}
        />
        <span className="text-sm text-text-secondary opacity-85">
          Words: {countWords(field.state.value as string)}
        </span>
      </div>
    );
  }

  return (
    <BaseTextField
      onChange={(val) => field.handleChange(val)}
      onBlur={field.handleBlur}
      {...props}
    />
  );
}
