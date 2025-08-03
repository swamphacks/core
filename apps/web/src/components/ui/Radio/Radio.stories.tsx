/* eslint-disable */
import { Form } from "react-aria-components";
import { Button } from "@/components/ui/Button";
import { Radio, RadioGroup } from ".";

export default {
  title: "UI/RadioGroup",
  component: RadioGroup,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {},
  args: {
    label: "How many hackathons have you participated in?",
    isDisabled: false,
    isRequired: false,
    description: "",
    children: (
      <>
        <Radio value="0">SwampHacks would be my first</Radio>
        <Radio value="1">1</Radio>
        <Radio value="2">2</Radio>
      </>
    ),
  },
};

export const Default = {
  args: {},
};

export const Validation = (args: any) => (
  <Form className="flex flex-col gap-2 items-start">
    <RadioGroup {...args} />
    <Button type="submit" variant="primary">
      Submit
    </Button>
  </Form>
);

Validation.args = {
  isRequired: true,
};
