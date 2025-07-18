import {
  TextField as RAC_TextField,
  type TextFieldProps as RAC_TextFieldProps,
  type ValidationResult,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import {
  Description,
  FieldError,
  Input,
  Label,
  fieldBorderStyles,
} from "@/components/ui/Field";
import { composeTailwindRenderProps, type Icon } from "@/components/ui/utils";

const inputStyles = tv({
  base: "outline-0 border-1 rounded-sm",
  variants: {
    isFocused: fieldBorderStyles.variants.isFocusWithin,
    isInvalid: fieldBorderStyles.variants.isInvalid,
    isDisabled: fieldBorderStyles.variants.isDisabled,
  },
});

export interface TextFieldProps extends RAC_TextFieldProps {
  label?: string;
  placeholder?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
  icon?: Icon;
}

const TextField = ({
  label,
  placeholder,
  description,
  errorMessage,
  isRequired,
  icon,
  ...props
}: TextFieldProps) => {
  return (
    <RAC_TextField
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "flex flex-col gap-1",
      )}
      isRequired={isRequired}
    >
      {label && <Label isRequired={isRequired}>{label}</Label>}
      <Input className={inputStyles} placeholder={placeholder} icon={icon} />

      {description && <Description>{description}</Description>}
      <FieldError>{errorMessage}</FieldError>
    </RAC_TextField>
  );
};

TextField.displayName = "TextField";

export { TextField };
