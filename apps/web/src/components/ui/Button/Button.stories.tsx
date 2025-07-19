import type { Meta, StoryObj } from "@storybook/react";
import { fn } from "@storybook/test";
import { Button } from ".";

export const ActionsData = {
  onClick: fn(),
};

const meta = {
  component: Button,
  title: "UI/Button",
  tags: ["autodocs"],
  args: {
    ...ActionsData,
    variant: "primary",
    children: "Button",
  },
  argTypes: {
    variant: {
      options: ["primary", "secondary", "danger"],
      control: { type: "select" },
    },
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    variant: "primary",
    children: "Primary Button",
  },
};

export const Secondary: Story = {
  args: {
    variant: "secondary",
    children: "Secondary Button",
  },
};

export const Danger: Story = {
  args: {
    variant: "danger",
    children: "Danger Button",
  },
};
