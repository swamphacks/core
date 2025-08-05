import React, { useEffect, useRef, useState } from "react";
import {
  Select as AriaSelect,
  type SelectProps as AriaSelectProps,
  Button,
  ListBox,
  type ListBoxItemProps,
  ListLayout,
  SelectValue,
  type ValidationResult,
  Virtualizer,
} from "react-aria-components";
import { tv } from "tailwind-variants";
import { Description, FieldError, Label } from "@/components/ui/Field";
import {
  DropdownItem,
  DropdownSection,
  type DropdownSectionProps,
} from "@/components/ui/ListBox";
import { Popover } from "@/components/ui/Popover";
import { composeTailwindRenderProps } from "@/components/ui/utils";
import TablerChevronDown from "~icons/tabler/chevron-down";
import { cn } from "@/utils/cn";

export const styles = tv({
  base: "h-9.5 flex items-center text-start gap-4 w-full cursor-default border border-input-border rounded-sm pl-3 pr-2 py-1.5 bg-surface",
  variants: {
    isDisabled: {
      false:
        "text-text-main hover:bg-gray-100 dark:hover:bg-neutral-700 group-invalid:border-red-600 forced-colors:group-invalid:border-[Mark]",
      true: "text-gray-200 dark:text-zinc-600 forced-colors:text-[GrayText] dark:bg-zinc-800 dark:border-white/5 forced-colors:border-[GrayText]",
    },
  },
});

export const listBoxContainerStyles = tv({
  base: "outline-hidden p-1 max-h-[inherit] overflow-auto [clip-path:inset(0_0_0_0_round_.75rem)]",
});

export interface SelectProps<T extends object>
  extends Omit<AriaSelectProps<T>, "children"> {
  label?: string;
  description?: string;
  errorMessage?: string | ((validation: ValidationResult) => string);
  items?: Iterable<T>;
  children: React.ReactNode | ((item: T) => React.ReactNode);
  virtualized?: boolean;
}

export function Select<T extends { id: string; name: string }>({
  label,
  description,
  errorMessage,
  items,
  virtualized = false,
  ...props
}: SelectProps<T>) {
  // This is needed in situation where the exact max width of the select is unknown and is controlled by something like `flex-1` for example
  const [maxWidth, setMaxWidth] = useState<number | undefined>(undefined);
  const selectRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (selectRef && selectRef.current) {
      setMaxWidth(selectRef.current.clientWidth);
    }
  }, []);

  return (
    <AriaSelect
      {...props}
      className={composeTailwindRenderProps(
        props.className,
        "group flex flex-col gap-1 font-figtree",
      )}
      style={{
        maxWidth,
      }}
      ref={selectRef}
    >
      {label && <Label isRequired={props.isRequired}>{label}</Label>}
      <Button className={styles}>
        <SelectValue className="flex-1 placeholder-shown:text-[#89898A] truncate overflow-hidden" />
        <TablerChevronDown
          aria-hidden
          className="w-4 h-4 text-gray-600 dark:text-zinc-400 forced-colors:text-[ButtonText] group-disabled:text-gray-200 dark:group-disabled:text-zinc-600 forced-colors:group-disabled:text-[GrayText]"
        />
      </Button>
      {description && <Description>{description}</Description>}
      <FieldError>{errorMessage}</FieldError>
      <Popover className="w-(--trigger-width)">
        {virtualized ? (
          <Virtualizer
            layout={ListLayout}
            layoutOptions={{ rowHeight: 33, padding: 4, gap: 0 }}
          >
            <ListBox
              items={items}
              className={cn(listBoxContainerStyles(), "p-0")}
            >
              {(item) => <SelectItem>{item.name}</SelectItem>}
            </ListBox>
          </Virtualizer>
        ) : (
          <ListBox items={items} className={listBoxContainerStyles()}>
            {(item) => <SelectItem>{item.name}</SelectItem>}
          </ListBox>
        )}
      </Popover>
    </AriaSelect>
  );
}

export function SelectItem(props: ListBoxItemProps) {
  return <DropdownItem {...props} />;
}

export function SelectSection<T extends object>(
  props: DropdownSectionProps<T>,
) {
  return <DropdownSection {...props} />;
}
