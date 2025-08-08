/* eslint-disable */
import type { Meta } from "@storybook/react";
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { ComboBox, ComboBoxItem, ComboBoxSection } from "./index";

const meta = {
  title: "UI/ComboBox",
  component: ComboBox,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
} satisfies Meta<typeof ComboBox>;

export default meta;

const items = [
  { id: "chocolate", name: "Chocolate" },
  { id: "mint", name: "Mint" },
  { id: "strawberry", name: "Strawberry" },
  { id: "vanilla", name: "Vanilla" },
];

export const Example = (args: any) => (
  <ComboBox
    itm
    label="Example"
    placeholder="Select an item"
    isRequired
    items={items}
    {...args}
  />
);

export const DisabledItems = (args: any) => <Example {...args} />;
DisabledItems.args = {
  disabledKeys: ["mint"],
};

export const Sections = (args: any) => (
  <ComboBox {...args}>
    <ComboBoxSection title="Fruit">
      <ComboBoxItem id="Apple">Apple</ComboBoxItem>
      <ComboBoxItem id="Banana">Banana</ComboBoxItem>
      <ComboBoxItem id="Orange">Orange</ComboBoxItem>
      <ComboBoxItem id="Honeydew">Honeydew</ComboBoxItem>
      <ComboBoxItem id="Grapes">Grapes</ComboBoxItem>
      <ComboBoxItem id="Watermelon">Watermelon</ComboBoxItem>
      <ComboBoxItem id="Cantaloupe">Cantaloupe</ComboBoxItem>
      <ComboBoxItem id="Pear">Pear</ComboBoxItem>
    </ComboBoxSection>

    <ComboBoxSection title="Another test">
      <ComboBoxItem id="Test1">Cabbage</ComboBoxItem>
      <ComboBoxItem id="Test2">Broccoli</ComboBoxItem>
      <ComboBoxItem id="Test3">Carrots</ComboBoxItem>
      <ComboBoxItem id="Test4">Lettuce</ComboBoxItem>
      <ComboBoxItem id="Test5">Spinach</ComboBoxItem>
      <ComboBoxItem id="Test7">Bok Choy</ComboBoxItem>
      <ComboBoxItem id="Test8">Cauliflower</ComboBoxItem>
      <ComboBoxItem id="Test9">Potatoes</ComboBoxItem>
    </ComboBoxSection>

    <ComboBoxSection title="Another test2">
      <ComboBoxItem id="Test11">Cabbage</ComboBoxItem>
      <ComboBoxItem id="Test21">Broccoli</ComboBoxItem>
      <ComboBoxItem id="Test31">Carrots</ComboBoxItem>
      <ComboBoxItem id="Test41">Lettuce</ComboBoxItem>
      <ComboBoxItem id="Test51">Spinach</ComboBoxItem>
      <ComboBoxItem id="Test71">Bok Choy</ComboBoxItem>
      <ComboBoxItem id="Test81">Cauliflower</ComboBoxItem>
      <ComboBoxItem id="Test91">Potatoes</ComboBoxItem>
    </ComboBoxSection>
  </ComboBox>
);

Sections.args = {
  label: "Preferred fruit or vegetable",
};

export const Validation = (args: any) => (
  <Form className="flex flex-col gap-2 items-start">
    <Example {...args} />
    <Button type="submit">Submit</Button>
  </Form>
);

Validation.args = {
  isRequired: true,
};
