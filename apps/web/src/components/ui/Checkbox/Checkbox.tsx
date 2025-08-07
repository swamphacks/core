// import { Check, Minus } from 'lucide-react';
import { type ReactNode } from "react";
import {
  Checkbox as AriaCheckbox,
  CheckboxGroup as AriaCheckboxGroup,
  type CheckboxGroupProps as AriaCheckboxGroupProps,
  type CheckboxProps,
  type ValidationResult,
  composeRenderProps,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import { Description, FieldError, Label } from "@/components/ui/Field";
import { composeTailwindRenderProps } from "@/components/ui/utils";
import TablerCheck from "~icons/tabler/check";
import TablerMinus from "~icons/tabler/minus";

export interface CheckboxGroupProps
  extends Omit<AriaCheckboxGroupProps, "children"> {
  label?: string;
  children?: ReactNode;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
}

export function CheckboxGroup({ label = "", ...props }: CheckboxGroupProps) {
  return (
    <AriaCheckboxGroup
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "flex flex-col gap-2",
      )}
    >
      <Label isRequired={props.isRequired}>{label}</Label>
      {props.children}
      {props.description && <Description>{props.description}</Description>}
      <FieldError>{props.errorMessage}</FieldError>
    </AriaCheckboxGroup>
  );
}

const checkboxStyles = tv({
  base: "flex gap-2 items-center group text-sm transition relative",
  variants: {
    isDisabled: {
      false: "text-gray-800 dark:text-zinc-200",
      true: "text-gray-300 dark:text-zinc-600 forced-colors:text-[GrayText]",
    },
  },
});

const boxStyles = tv({
  base: "w-5 h-5 shrink-0 rounded-sm flex items-center justify-center border-2 transition",
  variants: {
    isSelected: {
      false:
        "bg-white dark:bg-zinc-900 border-(--color) [--color:var(--color-gray-400)] dark:[--color:colors.zinc-400)] group-pressed:[--color:var(--color-gray-500)] dark:group-pressed:[--color:var(--color-zinc-300)]",
      true: "bg-(--color) border-(--color) [--color:var(--color-blue-600)] group-pressed:[--color:var(--color-blue-700)] dark:[--color:var(--color-blue-300)] dark:group-pressed:[--color:var(--color-blue-200)] forced-colors:[--color:Highlight]!",
    },
    isInvalid: {
      true: "[--color:var(--color-red-700)] dark:[--color:var(--color-red-600)] forced-colors:[--color:Mark]! group-pressed:[--color:var(--color-red-800)] dark:group-pressed:[--color:var(--color-red-700)]",
    },
    isDisabled: {
      true: "[--color:var(--color-gray-200)] dark:[--color:var(--color-zinc-700)] forced-colors:[--color:GrayText]!",
    },
  },
});

const iconStyles =
  "w-4 h-4 text-white group-disabled:text-gray-400 dark:text-slate-900 dark:group-disabled:text-slate-600 forced-colors:text-[HighlightText]";

export function Checkbox(props: CheckboxProps) {
  return (
    <AriaCheckbox
      {...props}
      className={composeRenderProps(props.className, (className, renderProps) =>
        checkboxStyles({ ...renderProps, className }),
      )}
    >
      {({ isSelected, isIndeterminate, ...renderProps }) => (
        <>
          <div
            className={boxStyles({
              isSelected: isSelected || isIndeterminate,
              ...renderProps,
            })}
          >
            {isIndeterminate ? (
              <TablerMinus aria-hidden className={iconStyles} />
            ) : isSelected ? (
              <TablerCheck aria-hidden className={iconStyles} />
            ) : null}
          </div>
          {props.children}
        </>
      )}
    </AriaCheckbox>
  );
}
