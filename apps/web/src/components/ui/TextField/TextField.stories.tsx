import type { Meta, StoryObj } from "@storybook/react";
import { TextField } from ".";
import TablerBrandLinkedin from "~icons/tabler/brand-linkedin";

const meta = {
  component: TextField,
  title: "UI/TextField",
  tags: ["autodocs"],
  args: {
    className: "max-w-70",
  },
} satisfies Meta<typeof TextField>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    label: "First Name",
    placeholder: "e.g. John",
    description: "This is a sample description of the text field.",
    isRequired: true,
  },
};

export const WithIcon: Story = {
  args: {
    label: "First Name",
    placeholder: "e.g. John",
    description: "This is a sample description of the text field.",
    isRequired: true,
    icon: TablerBrandLinkedin,
  },
};
