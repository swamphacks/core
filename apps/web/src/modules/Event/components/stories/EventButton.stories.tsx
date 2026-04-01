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
    .filter((key) => key !== "variant")
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
  title: "Event/Event Button",
  tags: ["autodocs"],
  argTypes,
} satisfies Meta<typeof EventButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Rejected: Story = {
  args: {
    status: "rejected",
    eventId: "example-event-id",
  },
};

export const Attending: Story = {
  args: {
    status: "attending",
    eventId: "example-event-id",
  },
};

export const Accepted: Story = {
  args: {
    status: "accepted",
    eventId: "example-event-id",
  },
};

export const Waitlisted: Story = {
  args: {
    status: "waitlisted",
    eventId: "example-event-id",
  },
};

export const UnderReview: Story = {
  args: {
    status: "underReview",
    eventId: "example-event-id",
  },
};

export const NotApplied: Story = {
  args: {
    status: "notApplied",
    eventId: "example-event-id",
  },
};

export const Staff: Story = {
  args: {
    status: "staff",
    eventId: "example-event-id",
  },
};

export const Admin: Story = {
  args: {
    status: "admin",
    eventId: "example-event-id",
  },
};

export const NotGoing: Story = {
  args: {
    status: "notGoing",
    eventId: "example-event-id",
  },
};

export const Completed: Story = {
  args: {
    status: "completed",
    eventId: "example-event-id",
  },
};
