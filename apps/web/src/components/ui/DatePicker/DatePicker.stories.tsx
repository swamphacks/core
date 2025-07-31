import type { Meta } from "@storybook/react";
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { DatePicker } from ".";

const meta: Meta<typeof DatePicker> = {
  title: "UI/DatePicker",
  component: DatePicker,
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
export const Example = (args: any) => <DatePicker {...args} />;

/* eslint-disable @typescript-eslint/no-explicit-any */
export const Validation = (args: any) => (
  <Form className="flex flex-col gap-2 items-start">
    <DatePicker {...args} />
    <Button type="submit" variant="secondary">
      Submit
    </Button>
  </Form>
);

Validation.args = {
  isRequired: true,
};
