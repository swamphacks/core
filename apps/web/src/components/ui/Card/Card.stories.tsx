import type { Meta, StoryObj } from "@storybook/react";
import { Card } from ".";

const meta = {
  component: Card,
  title: "UI/Card",
  tags: ["autodocs"],
} satisfies Meta<typeof Card>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: (
      <div>
        <p>Title</p>
        <p>This is the card's description</p>
      </div>
    ),
    size: "sm",
  },
};
