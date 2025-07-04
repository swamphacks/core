import type { Meta, StoryObj } from "@storybook/react";
import { EventBadge, eventBadge } from "../EventBadge";

const entries = [];

for (let i = 0; i < eventBadge.variantKeys.length; i++) {
  const key = eventBadge.variantKeys[i];

  if (key === "type") continue;

  entries.push([
    key,
    {
      control: {
        type: "select",
      },
      options: Object.keys(eventBadge.variants[key]),
    },
  ]);
}
const argTypes = Object.fromEntries(entries);

const meta = {
  component: EventBadge,
  title: "UI/Event Badge",
  tags: ["autodocs"],
  argTypes,
} satisfies Meta<typeof EventBadge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Rejected: Story = {
  args: {
    status: "rejected",
    size: "sm",
  },
};

export const Attending: Story = {
  args: {
    status: "attending",
    size: "sm",
  },
};

export const Accepted: Story = {
  args: {
    status: "accepted",
    size: "sm",
  },
};

export const Waitlisted: Story = {
  args: {
    status: "waitlisted",
    size: "sm",
  },
};

export const UnderReview: Story = {
  args: {
    status: "underReview",
    size: "sm",
  },
};

export const NotApplied: Story = {
  args: {
    status: "notApplied",
    size: "sm",
  },
};

export const Staff: Story = {
  args: {
    status: "staff",
    size: "sm",
  },
};

export const Admin: Story = {
  args: {
    status: "admin",
    size: "sm",
  },
};

export const NotGoing: Story = {
  args: {
    status: "notGoing",
    size: "sm",
  },
};

export const Completed: Story = {
  args: {
    status: "completed",
    size: "sm",
  },
};
