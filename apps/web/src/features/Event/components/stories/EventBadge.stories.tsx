import type { Meta, StoryObj } from "@storybook/react";
import { EventBadge } from "../EventBadge";
import applicationStatus from "../../applicationStatus";

const meta = {
  component: EventBadge,
  title: "UI/Event Badge",
  tags: ["autodocs"],
  argTypes: {
    status: {
      control: {
        type: "select",
      },
      options: Object.keys(applicationStatus),
    },
    size: {
      control: {
        type: "select",
      },
      // Maybe we can figure out a way to extract the sizes from tailwind variants?
      options: ["sm", "md"],
    },
  },
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
