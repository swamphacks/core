import {
  type FieldErrorProps,
  Group,
  type GroupProps,
  type InputProps,
  type LabelProps,
  FieldError as RACFieldError,
  Input as RACInput,
  Label as RACLabel,
  Text,
  type TextProps,
  composeRenderProps,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import { composeTailwindRenderProps, type Icon } from "@/components/ui/utils";
import { cn } from "@/utils/cn";
import TablerAsterisk from "~icons/tabler/asterisk";
import { forwardRef } from "react";

export const fieldBorderStyles = tv({
  variants: {
    isFocusWithin: {
      false: "border-input-border forced-colors:border-[ButtonBorder]",
      true: "border-input-border-focused forced-colors:border-[Highlight]",
    },
    isInvalid: {
      true: "border-input-border-invalid forced-colors:border-[Mark]",
    },
    isDisabled: {
      true: "border-input-border-disabled forced-colors:border-[GrayText]",
    },
  },
});

export const fieldGroupStyles = tv({
  base: "group flex items-start bg-surfaceforced-colors:bg-[Field]",
  variants: fieldBorderStyles.variants,
});

export function Label({
  isRequired,
  ...props
}: LabelProps & { isRequired?: boolean }) {
  return (
    <RACLabel
      {...props}
      className={cn(
        "flex items-center gap-1 text-text-secondary font-medium cursor-default w-fit",
        props.className,
      )}
    >
      {props.children}
      {isRequired && <TablerAsterisk className="text-[8px] text-red-500" />}
    </RACLabel>
  );
}

export function Description(props: TextProps) {
  return (
    <Text
      {...props}
      slot="description"
      className={cn("text-sm text-text-secondary opacity-85", props.className)}
    />
  );
}

export function ErrorList({ errors }: { errors: string[] }) {
  return (
    <ul className="list-disc ml-5">
      {errors.map((error, i) => (
        <li key={i}>{error}</li>
      ))}
    </ul>
  );
}

export function FieldError(props: FieldErrorProps) {
  return (
    <RACFieldError
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "text-sm text-input-text-error forced-colors:text-[Mark]",
      )}
    >
      {props.children
        ? props.children
        : ({ validationErrors }) => {
            return validationErrors.length === 1 ? (
              validationErrors[0]
            ) : (
              <ErrorList errors={validationErrors} />
            );
          }}
    </RACFieldError>
  );
}

export function FieldGroup(props: GroupProps) {
  return (
    <Group
      {...props}
      className={composeRenderProps(props.className, (className, renderProps) =>
        fieldGroupStyles({ ...renderProps, className }),
      )}
    />
  );
}

export const Input = forwardRef(
  (
    { icon: Icon, ...props }: InputProps & { icon?: Icon },
    ref: React.ForwardedRef<HTMLInputElement>,
  ) => {
    const inputClassName =
      "h-9 py-1.5 w-full min-w-0 outline-0 bg-surface text-base text-text-main disabled:text-input-text-disabled";

    if (!Icon) {
      return (
        <RACInput
          {...props}
          className={composeTailwindRenderProps(
            props.className,
            cn(inputClassName, "px-2"),
          )}
          ref={ref}
        />
      );
    }

    return (
      <div className="flex items-center gap-1 bg-surface relative">
        <div className="ml-2 text-text-secondary absolute pointer-events-none opacity-85">
          <Icon />
        </div>

        <RACInput
          {...props}
          className={composeTailwindRenderProps(
            props.className,
            cn(inputClassName, "pl-8 pr-2"),
          )}
          ref={ref}
        />
      </div>
    );
  },
);
