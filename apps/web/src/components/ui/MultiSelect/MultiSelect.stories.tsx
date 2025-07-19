/* eslint-disable */
import type { Meta } from "@storybook/react";
import { MultiSelect } from ".";

const meta: Meta<typeof MultiSelect> = {
  title: "UI/MultiSelect",
  component: MultiSelect,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  args: {
    name: "Ice cream flavor",
  },
};

export default meta;

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

export const Example = (args: any) => (
  <MultiSelect
    name="majors"
    label="Major(s)"
    options={options}
    isRequired
    {...args}
  />
);
