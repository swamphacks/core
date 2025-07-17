/* eslint-disable */
import {
  type AriaPopoverProps,
  type AriaSelectProps,
  HiddenSelect,
  type Key,
  Overlay,
  usePopover,
  useSelect,
} from "react-aria";
import { DismissButton } from "react-aria";
import React, { type JSX } from "react";
import { useButton } from "react-aria";
import { useListBox, useOption } from "react-aria";
import {
  useSelectState,
  Item,
  type OverlayTriggerState,
  useListState,
} from "react-stately";

function Select<T extends object>(props: AriaSelectProps<T>): JSX.Element {
  // Create state based on the incoming props
  let state = useSelectState(props);

  // Get props for child elements from useSelect
  let ref = React.useRef(null);
  let { labelProps, triggerProps, valueProps, menuProps } = useSelect(
    props,
    state,
    ref,
  );

  let [selected, setSelected] = React.useState<Set<Key>>(new Set());

  const onSelectionChange = (selection: any) => {
    setSelected(selection);
  };

  return (
    <div style={{ display: "inline-block" }}>
      <div {...labelProps}>{props.label}</div>
      <HiddenSelect
        isDisabled={props.isDisabled}
        state={state}
        triggerRef={ref}
        label={props.label}
        name={props.name}
      />
      <Button
        {...triggerProps}
        buttonRef={ref}
        style={{ height: 30, fontSize: 14 }}
      >
        <span {...valueProps}>
          {selected.size === 0 ? "select an item" : [...selected].join(",")}
        </span>
        <span aria-hidden="true" style={{ paddingLeft: 5 }}>
          ▼
        </span>
      </Button>
      {state.isOpen && (
        <Popover state={state} triggerRef={ref} placement="bottom start">
          <ListBox
            {...menuProps}
            state={state}
            onSelectionChange={onSelectionChange}
            selected={selected}
          />
        </Popover>
      )}
    </div>
  );
}

interface PopoverProps extends Omit<AriaPopoverProps, "popoverRef"> {
  children: React.ReactNode;
  state: OverlayTriggerState;
}

function Popover({ children, state, ...props }: PopoverProps) {
  let popoverRef = React.useRef(null);
  let { popoverProps, underlayProps } = usePopover(
    {
      ...props,
      popoverRef,
    },
    state,
  );

  return (
    <Overlay>
      <div {...underlayProps} style={{ position: "fixed", inset: 0 }} />
      <div
        {...popoverProps}
        ref={popoverRef}
        style={{
          ...popoverProps.style,
          background: "var(--page-background)",
          border: "1px solid gray",
        }}
      >
        <DismissButton onDismiss={state.close} />
        {children}
        <DismissButton onDismiss={state.close} />
      </div>
    </Overlay>
  );
}

// @ts-ignore
function ListBox(props) {
  let ref = React.useRef(null);
  let { listBoxRef = ref, state: _state } = props;

  const state = useListState({
    collection: _state.collection,
    disabledKeys: [],
    selectionMode: "multiple",
    onSelectionChange: (selection) => props.onSelectionChange(selection),
    selectedKeys: props.selected,
  });
  let { listBoxProps } = useListBox(props, state, listBoxRef);

  return (
    <ul
      {...listBoxProps}
      ref={listBoxRef}
      style={{
        margin: 0,
        padding: 0,
        listStyle: "none",
        maxHeight: 150,
        overflow: "auto",
        minWidth: 100,
        background: "lightgray",
      }}
    >
      {[...state.collection].map((item) => (
        <Option key={item.key} item={item} state={state} />
      ))}
    </ul>
  );
}

// @ts-ignore
function Option({ item, state }) {
  let ref = React.useRef(null);
  let { optionProps, isSelected, isFocused, isDisabled } = useOption(
    { key: item.key },
    state,
    ref,
  );

  return (
    <li
      {...optionProps}
      ref={ref}
      style={{
        background: isFocused ? "gray" : "transparent",
        color: isDisabled ? "gray" : isFocused ? "white" : "black",
        padding: "2px 5px",
        outline: "none",
        cursor: "pointer",
        display: "flex",
        justifyContent: "space-between",
        gap: "10px",
      }}
    >
      {item.rendered}
      {isSelected ? <span>✓</span> : null}
    </li>
  );
}

// @ts-ignore
function Button(props) {
  let ref = props.buttonRef;
  let { buttonProps } = useButton(props, ref);
  return (
    <button {...buttonProps} ref={ref} style={props.style}>
      {props.children}
    </button>
  );
}

let lotsOfItems: any[] = [];
for (let i = 0; i < 50; i++) {
  lotsOfItems.push({ name: "Item " + i });
}

function MultiSelect() {
  return (
    <Select label="Favorite Color">
      <Item key="red">Red</Item>
      <Item key="orange">Orange</Item>
      <Item key="yellow">Yellow</Item>
      <Item key="green">Green</Item>
      {/* <Item>Blue</Item>
      <Item>Purple</Item>
      <Item>Black</Item>
      <Item>White</Item>
      <Item>Lime</Item>
      <Item>Fushsia</Item> */}
    </Select>
  );
}

export { MultiSelect };
