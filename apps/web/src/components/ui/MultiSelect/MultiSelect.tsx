import Select, {
  components,
  type ClearIndicatorProps,
  type DropdownIndicatorProps,
  type MultiValueRemoveProps,
} from "react-select";
import TablerChevronDown from "~icons/tabler/chevron-down";
import {
  listBoxContainerStyles,
  styles as selectStyles,
} from "@/components/ui/Select";
import { styles as PopoverStyles } from "@/components/ui/Popover";
import { dropdownItemStyles } from "@/components/ui/ListBox";
import { cn } from "@/utils/cn";
import {
  removeButtonStyles as removeTagButtonStyles,
  tagStyles,
} from "@/components/ui/Tag";
import TablerX from "~icons/tabler/x";
import { useId, useRef, useState, type CSSProperties } from "react";
import { ErrorList, fieldBorderStyles } from "@/components/ui/Field";
import TablerAsterisk from "~icons/tabler/asterisk";

const DropdownIndicator = (props: DropdownIndicatorProps) => {
  return (
    <components.DropdownIndicator {...props}>
      <TablerChevronDown
        aria-hidden
        className="w-4 h-4 text-gray-600 dark:text-zinc-400 forced-colors:text-[ButtonText] group-disabled:text-gray-200 dark:group-disabled:text-zinc-600 forced-colors:group-disabled:text-[GrayText]"
      />
    </components.DropdownIndicator>
  );
};

const ClearIndicator = (props: ClearIndicatorProps) => {
  const {
    getStyles,
    innerProps: { ref, ...restInnerProps },
  } = props;
  return (
    <div
      {...restInnerProps}
      ref={ref}
      style={getStyles("clearIndicator", props) as CSSProperties}
      className="mr-2"
    >
      <TablerX />
    </div>
  );
};

const MultiValueRemove = (props: MultiValueRemoveProps) => {
  return (
    <components.MultiValueRemove {...props}>
      <TablerX aria-hidden className="size-3" />
    </components.MultiValueRemove>
  );
};

interface Option {
  value: string;
  label: string;
}

export interface MultiSelectProps {
  // Name must be unique and is required for form submissions
  name: string;
  label: string;
  options: Option[];
  isRequired?: boolean;
  onChange?: (data: Option[]) => void;

  // Not sure if its better to do this https://react-spectrum.adobe.com/react-aria/Form.html#custom-children
  // accepting errors props should work for now.
  errors?: string[];
}

export const MULTISELECT_NAME_PREFIX = "multiselect-";

const MultiSelect = ({
  isRequired,
  label,
  name,
  options,
  onChange,
  errors,
  ...props
}: MultiSelectProps) => {
  const id = useId();
  const hiddenSelectRef = useRef<HTMLSelectElement | null>(null);
  const [isInvalid, setIsInvalid] = useState(false);

  const renderErrors = () => {
    if (isInvalid) return <p>Please select an item in the list.</p>;

    if (errors && errors.length > 0) {
      return <ErrorList errors={errors} />;
    }
  };

  return (
    <div className="flex flex-col gap-1 flex-1 font-figtree">
      <label
        className="flex items-center gap-1 text-text-secondary font-medium cursor-default w-fit"
        htmlFor={id}
      >
        {label}
        {isRequired && <TablerAsterisk className="text-[8px] text-red-500" />}
      </label>

      <Select
        inputId={id}
        menuPlacement="auto"
        unstyled
        isMulti
        options={options}
        components={{ DropdownIndicator, ClearIndicator, MultiValueRemove }}
        classNames={{
          placeholder: () => "text-[#89898A]",
          control: (state) =>
            cn(
              selectStyles(state),
              "h-full",
              fieldBorderStyles({
                isInvalid: isInvalid || (errors && errors.length > 0),
              }),
            ),
          menu: () => cn(PopoverStyles(), "mt-2"),
          menuList: (state) => listBoxContainerStyles(state),
          option: (state) => dropdownItemStyles(state),
          multiValue: () => cn(tagStyles({ color: "gray" })),
          multiValueRemove: () => removeTagButtonStyles(),
          valueContainer: () => "gap-1",
        }}
        onChange={(data) => {
          const selectedValues = data.map((item) => (item as Option).value);

          if (hiddenSelectRef.current) {
            for (const option of hiddenSelectRef.current.options) {
              option.selected = selectedValues.includes(option.value);
            }
          }

          if (isInvalid && data.length > 0) {
            setIsInvalid(false);
          }

          onChange?.(data as Option[]);
        }}
        styles={{
          control: (baseStyles) => ({
            ...baseStyles,
            transition: "none",
          }),
        }}
        {...props}
      />

      <div className="text-sm text-input-text-error forced-colors:text-[Mark]">
        {renderErrors()}
      </div>

      {/* This is a dummy select element used for form submissions. Its values are updated based on the actual Select component above.
          It doesn't do anything since we are using Tanstack Form, but i'll leave it here just in case it comes in handy in the future.
      */}
      <select
        tabIndex={-1}
        autoComplete="off"
        style={{
          opacity: 0,
          height: 0,
          position: "absolute",
          pointerEvents: "none",
        }}
        required={isRequired}
        multiple
        name={`${MULTISELECT_NAME_PREFIX}${name}`}
        ref={hiddenSelectRef}
        onInvalid={() => setIsInvalid(true)}
      >
        {options.map((option) => (
          <option
            value={option.value}
            key={option.value}
            defaultValue={option.value}
          />
        ))}
      </select>
    </div>
  );
};

export { MultiSelect };
