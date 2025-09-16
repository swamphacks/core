import {
  Menu as AriaMenu,
  MenuItem as AriaMenuItem,
  type MenuProps as AriaMenuProps,
  type MenuItemProps,
  MenuSection as AriaMenuSection,
  type MenuSectionProps as AriaMenuSectionProps,
  Separator,
  type SeparatorProps,
  composeRenderProps,
  Header,
  Collection,
} from "react-aria-components";
import { dropdownItemStyles } from "../ListBox";
import { Popover, type PopoverProps } from "../Popover";
import TablerChevronRight from "~icons/tabler/chevron-right";
import TablerCheck from "~icons/tabler/check";
import { composeTailwindRenderProps } from "@/components/ui/utils";

interface MenuProps<T> extends AriaMenuProps<T> {
  placement?: PopoverProps["placement"];
  header?: React.ReactNode;
}

export function Menu<T extends object>({ header, ...props }: MenuProps<T>) {
  return (
    <Popover placement={props.placement} className="min-w-[170px]">
      {header}
      <AriaMenu
        {...props}
        className="p-1 outline max-h-[inherit] overflow-auto [clip-path:inset(0_0_0_0_round_.75rem)]"
      />
    </Popover>
  );
}

export function MenuItem({ className, ...props }: MenuItemProps) {
  const textValue =
    props.textValue ||
    (typeof props.children === "string" ? props.children : undefined);
  return (
    <AriaMenuItem
      textValue={textValue}
      {...props}
      className={composeTailwindRenderProps(
        dropdownItemStyles,
        className as string,
      )}
    >
      {composeRenderProps(
        props.children,
        (children, { selectionMode, isSelected, hasSubmenu }) => (
          <>
            {selectionMode !== "none" && (
              <span className="flex items-center w-4">
                {isSelected && <TablerCheck aria-hidden className="w-4 h-4" />}
              </span>
            )}
            <span className="flex items-center flex-1 gap-2 font-normal truncate group-selected:font-semibold">
              {children}
            </span>
            {hasSubmenu && (
              <TablerChevronRight
                aria-hidden
                className="absolute w-4 h-4 right-2"
              />
            )}
          </>
        ),
      )}
    </AriaMenuItem>
  );
}

export function MenuSeparator(props: SeparatorProps) {
  return (
    <Separator
      {...props}
      className="mx-3 my-1 border-b border-gray-300 dark:border-zinc-700"
    />
  );
}

export interface MenuSectionProps<T> extends AriaMenuSectionProps<T> {
  title?: string;
  items?: any;
}

export function MenuSection<T extends object>(props: MenuSectionProps<T>) {
  return (
    <AriaMenuSection
      {...props}
      className="first:-mt-[5px] after:content-[''] after:block after:h-[5px]"
    >
      <Header className="text-sm font-semibold text-gray-500 dark:text-zinc-300 px-4 py-1 truncate sticky -top-[5px] -mt-px -mx-1 z-10 bg-gray-100/60 dark:bg-zinc-700/60 backdrop-blur-md supports-[-moz-appearance:none]:bg-gray-100 border-y border-y-gray-200 dark:border-y-zinc-700 [&+*]:mt-1">
        {props.title}
      </Header>
      <Collection items={props.items}>{props.children}</Collection>
    </AriaMenuSection>
  );
}
