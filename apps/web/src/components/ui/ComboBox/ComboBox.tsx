import React from "react";
import {
  ComboBox as AriaComboBox,
  type ComboBoxProps as AriaComboBoxProps,
  ListBox,
  type ListBoxItemProps,
  type ValidationResult,
  Virtualizer,
  ListLayout,
} from "react-aria-components";
import { Button } from "@/components/ui/Button";
import {
  Description,
  FieldError,
  FieldGroup,
  Input,
  Label,
} from "@/components/ui/Field";
import {
  DropdownItem,
  DropdownSection,
  type DropdownSectionProps,
} from "@/components/ui/ListBox";
import { Popover } from "@/components/ui/Popover";
import { composeTailwindRenderProps } from "@/components/ui/utils";
import TablerChevronDown from "~icons/tabler/chevron-down";

export interface ComboBoxProps<T extends { id: string; name: string }>
  extends Omit<AriaComboBoxProps<T>, "children"> {
  label?: string;
  description?: string | null;
  errorMessage?: string | ((validation: ValidationResult) => string);
  children: React.ReactNode | ((item: T) => React.ReactNode);
  placeholder?: string;
  virtualized?: boolean;
}

export function ComboBox<T extends { id: string; name: string }>({
  label,
  description,
  errorMessage,
  items,
  placeholder,
  children,
  virtualized = false,
  ...props
}: ComboBoxProps<T>) {
  return (
    <AriaComboBox
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group flex flex-col gap-1 font-figtree",
      )}
    >
      <Label isRequired={props.isRequired}>{label}</Label>
      <FieldGroup className="border-1 rounded-sm bg-surface overflow-hidden items-center">
        <Input placeholder={placeholder} />
        <Button
          variant="icon"
          size="auto"
          className="w-6 mr-1 rounded-sm outline-offset-0"
        >
          <TablerChevronDown aria-hidden className="w-4 h-4" />
        </Button>
      </FieldGroup>
      {description && <Description>{description}</Description>}
      <FieldError>{errorMessage}</FieldError>
      <Popover
        className="min-w-[200px] sm:min-w-(--trigger-width) ml-0.5"
        isNonModal
      >
        {virtualized ? (
          <Virtualizer
            layout={ListLayout}
            layoutOptions={{ rowHeight: 33, padding: 4, gap: 2 }}
          >
            <ListBox
              items={items}
              className="overflow-y-auto max-h-70 outline-0 p-0 w-full box-border m-0"
            >
              {children
                ? children
                : (item) => <ComboBoxItem>{item.name}</ComboBoxItem>}
            </ListBox>
          </Virtualizer>
        ) : (
          <ListBox
            items={items}
            className="outline-0 p-1 overflow-auto max-h-70"
          >
            {children
              ? children
              : (item) => <ComboBoxItem>{item.name}</ComboBoxItem>}
          </ListBox>
        )}
      </Popover>
    </AriaComboBox>
  );
}

export function ComboBoxItem(props: ListBoxItemProps) {
  return <DropdownItem {...props} />;
}

export function ComboBoxSection<T extends object>(
  props: DropdownSectionProps<T>,
) {
  return <DropdownSection {...props} />;
}
