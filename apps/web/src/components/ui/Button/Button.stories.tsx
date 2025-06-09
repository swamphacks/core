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
    color: "primary",
    children: "Button",
  },
  argTypes: {
    color: {
      options: ["primary", "secondary", "danger"],
      control: { type: "select" },
    },
  },
} satisfies Meta<typeof Button>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: {
    color: "primary",
    children: "Primary Button",
  },
};

export const Secondary: Story = {
  args: {
    color: "secondary",
    children: "Secondary Button",
  },
};

export const Danger: Story = {
  args: {
    color: "danger",
    children: "Danger Button",
  },
};
