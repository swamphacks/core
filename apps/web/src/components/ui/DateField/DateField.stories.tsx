import type { Meta } from "@storybook/react";
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { DateField } from ".";

const meta: Meta<typeof DateField> = {
  title: "UI/DateField",
  component: DateField,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  args: {
    label: "Event date",
  },
};

export default meta;

/* eslint-disable @typescript-eslint/no-explicit-any */
export const Example = (args: any) => <DateField {...args} />;

/* eslint-disable @typescript-eslint/no-explicit-any */
export const Validation = (args: any) => (
  <Form className="flex flex-col gap-2 items-start">
    <DateField {...args} />
    <Button type="submit" variant="secondary">
      Submit
    </Button>
  </Form>
);

Validation.args = {
  isRequired: true,
};
