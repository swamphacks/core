import type { Meta, StoryObj } from "@storybook/react";
import { EventButton, eventButton } from "../EventButton";

const entries = [];

for (let i = 0; i < eventButton.variantKeys.length; i++) {
  const key = eventButton.variantKeys[i];

  if (key === "color") continue;

  entries.push([
    key,
    {
      control: {
        type: "select",
      },
      options: Object.keys(eventButton.variants[key]),
    },
  ]);
}
const argTypes = Object.fromEntries(entries);

const meta = {
  component: EventButton,
  title: "UI/Event Button",
  tags: ["autodocs"],
  argTypes,
} satisfies Meta<typeof EventButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Rejected: Story = {
  args: {
    status: "rejected",
  },
};

export const Attending: Story = {
  args: {
    status: "attending",
  },
};

export const Accepted: Story = {
  args: {
    status: "accepted",
  },
};

export const Waitlisted: Story = {
  args: {
    status: "waitlisted",
  },
};

export const UnderReview: Story = {
  args: {
    status: "underReview",
  },
};

export const NotApplied: Story = {
  args: {
    status: "notApplied",
  },
};

export const Staff: Story = {
  args: {
    status: "staff",
  },
};

export const Admin: Story = {
  args: {
    status: "admin",
  },
};

export const NotGoing: Story = {
  args: {
    status: "notGoing",
  },
};

export const Completed: Story = {
  args: {
    status: "completed",
  },
};
