import { type ReactNode } from "react";
import {
  RadioGroup as RACRadioGroup,
  type RadioGroupProps as RACRadioGroupProps,
  type ValidationResult,
} from "react-aria-components";
import { Description, FieldError, Label } from "@/components/ui/Field";
import { composeTailwindRenderProps } from "@/components/ui/utils";

export interface RadioGroupProps extends Omit<RACRadioGroupProps, "children"> {
  label?: string;
  children?: ReactNode;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
}

export function RadioGroup(props: RadioGroupProps) {
  return (
    <RACRadioGroup
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group flex flex-col gap-2 font-figtree",
      )}
    >
      <Label isRequired={props.isRequired}>{props.label}</Label>
      <div className="flex group-orientation-vertical:flex-col gap-2 group-orientation-horizontal:gap-4">
        {props.children}
      </div>
      {props.description && <Description>{props.description}</Description>}
      <FieldError>{props.errorMessage}</FieldError>
    </RACRadioGroup>
  );
}
