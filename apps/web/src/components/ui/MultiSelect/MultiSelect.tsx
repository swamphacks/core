/* eslint-disable */
import Select, { components, type DropdownIndicatorProps } from "react-select";
import TablerChevronDown from "~icons/tabler/chevron-down";
import { listBoxContainerStyles, styles as selectStyles } from "../Select";
import { styles as PopoverStyles } from "../Popover";
import { dropdownItemStyles } from "../ListBox";
import { cn } from "@/utils/cn";
import { removeButtonStyles as removeTagButtonStyles, tagStyles } from "../Tag";
import TablerX from "~icons/tabler/x";
import { useId } from "react";

const options = [
  { value: "chocolate", label: "Chocolate" },
  { value: "strawberry", label: "Strawberry" },
  { value: "vanilla", label: "Vanilla" },
  { value: "vanilla2", label: "Vanilla" },
  { value: "vanilla3", label: "Vanilla" },
  { value: "vanilla4", label: "Vanilla" },
  { value: "vanilla5", label: "Vanilla" },
  { value: "vanilla6", label: "Vanilla" },
  { value: "vanilla7", label: "Vanilla" },
  { value: "vanilla8", label: "Vanilla" },
];

interface MultiSelectProps {
  isRequired?: boolean;
  name: string;
}

const DropdownIndicator = (
  props: DropdownIndicatorProps<(typeof options)[0], true>,
) => {
  return (
    <components.DropdownIndicator {...props}>
      <TablerChevronDown
        aria-hidden
        className="w-4 h-4 text-gray-600 dark:text-zinc-400 forced-colors:text-[ButtonText] group-disabled:text-gray-200 dark:group-disabled:text-zinc-600 forced-colors:group-disabled:text-[GrayText]"
      />
    </components.DropdownIndicator>
  );
};

// @ts-ignore
const ClearIndicator = (props) => {
  const {
    getStyles,
    innerProps: { ref, ...restInnerProps },
  } = props;
  return (
    <div
      {...restInnerProps}
      ref={ref}
      style={getStyles("clearIndicator", props)}
      className="mr-2"
    >
      <TablerX />
    </div>
  );
};

// @ts-ignore
const MultiValueRemove = (props) => {
  return (
    <components.MultiValueRemove {...props}>
      <TablerX aria-hidden className="size-3" />
    </components.MultiValueRemove>
  );
};

const MultiSelect = ({ isRequired, ...props }: MultiSelectProps) => {
  const id = useId();

  return (
    <div className="flex flex-col gap-1 flex-1 font-figtree">
      <label htmlFor={id}>{props.name}</label>
      <Select
        inputId={id}
        // aria-labelledby={id}
        unstyled
        isMulti
        options={options}
        components={{ DropdownIndicator, ClearIndicator, MultiValueRemove }}
        classNames={{
          placeholder: () => "text-[#89898A]",
          control: (state) => cn(selectStyles(state), "h-full"),
          menu: () => cn(PopoverStyles(), "mt-2"),
          menuList: (state) => listBoxContainerStyles(state),
          option: (state) => dropdownItemStyles(state),
          multiValue: () => cn(tagStyles({ color: "gray" })),
          multiValueRemove: () => removeTagButtonStyles(),
          valueContainer: () => "gap-1",
        }}
        styles={{
          control: (baseStyles) => ({
            ...baseStyles,
            transition: "none",
          }),
        }}
        {...props}
      />
    </div>
  );
};

export { MultiSelect };
