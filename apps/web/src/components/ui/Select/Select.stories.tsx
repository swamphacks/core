/* eslint-disable */
import type { Meta } from "@storybook/react";
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { Select, SelectItem, SelectSection } from ".";

const meta: Meta<typeof Select> = {
  title: "UI/Select",
  component: Select,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  args: {
    label: "Ice cream flavor",
  },
};

export default meta;

const data = [
  {
    id: "Small",
    name: "Small",
  },
  {
    id: "Medium",
    name: "Medium",
  },
  {
    id: "Large",
    name: "Large",
  },
];

export const Example = (args: any) => (
  <Select
    {...args}
    label="T-Shirt Size"
    placeholder="Select size"
    items={data}
  />
);

export const DisabledItems = (args: any) => <Example {...args} />;
DisabledItems.args = {
  disabledKeys: ["Small"],
};

export const Validation = (args: any) => (
  <Form className="flex flex-col gap-2 items-start">
    <Example {...args} />
    <Button type="submit" variant="secondary">
      Submit
    </Button>
  </Form>
);

Validation.args = {
  isRequired: true,
};
