import type { Meta, StoryObj } from "@storybook/react";
import { EventButton, eventButton } from "../EventButton";

/**
 * Generates the Storybook `argTypes` configuration dynamically based on
 * the available variants in `eventButton`.
 *
 * Filter out the `color` variant since it is not relevant for the storybook options.
 */
const argTypes = Object.fromEntries(
  eventButton.variantKeys
    .filter((key) => key !== "color")
    .map((key) => [
      key,
      {
        control: {
          type: "select",
        },
        options: Object.keys(eventButton.variants[key]),
      },
    ]),
);

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
