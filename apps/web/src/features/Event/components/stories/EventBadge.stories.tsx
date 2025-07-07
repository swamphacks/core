import type { Meta, StoryObj } from "@storybook/react";
import { EventBadge, eventBadge } from "../EventBadge";

/**
 * Generates the Storybook `argTypes` configuration dynamically based on
 * the available variants in `eventBadge`.
 *
 * Filter out the `type` variant since it is not relevant for the storybook options.
 */
const argTypes = Object.fromEntries(
  eventBadge.variantKeys
    .filter((key) => key !== "type")
    .map((key) => [
      key,
      {
        control: {
          type: "select",
        },
        options: Object.keys(eventBadge.variants[key]),
      },
    ]),
);

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
