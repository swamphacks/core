import { Checkbox } from ".";

export default {
  title: "UI/Checkbox",
  component: Checkbox,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
  argTypes: {},
  args: {
    isDisabled: false,
    children: "Checkbox",
  },
};

export const Default = {
  args: {},
};
