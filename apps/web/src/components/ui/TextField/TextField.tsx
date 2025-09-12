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
import { TextArea } from "react-aria-components";
import { cn } from "@/utils/cn";

export const inputStyles = tv({
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
  iconPlacement?: "left" | "right";
  textarea?: boolean;
}

const TextField = ({
  label,
  placeholder,
  description,
  errorMessage,
  isRequired,
  icon,
  textarea,
  iconPlacement = "left",
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
      {textarea ? (
        <TextArea
          className={composeTailwindRenderProps(
            inputStyles,
            cn(
              "px-2 py-1.5 bg-input-bg text-base text-text-main disabled:cursor-not-allowed disabled:text-input-text-disabled disabled:bg-input-bg-disbaled",
            ),
          )}
          placeholder={placeholder}
        />
      ) : (
        <Input
          className={inputStyles}
          placeholder={placeholder}
          icon={icon}
          iconPlacement={iconPlacement}
        />
      )}

      {description && <Description>{description}</Description>}
      <FieldError>{errorMessage}</FieldError>
    </RAC_TextField>
  );
};

TextField.displayName = "TextField";

export { TextField };
