import type { Meta, StoryObj } from "@storybook/react";
import { EventCard as EventCardComp } from "../EventCard";
import applicationStatus from "../../applicationStatus";

const meta = {
  component: EventCardComp,
  title: "UI/Event Card",
  tags: ["autodocs"],
  argTypes: {
    status: {
      control: {
        type: "select",
      },
      options: Object.keys(applicationStatus),
    },
  },
} satisfies Meta<typeof EventCardComp>;

export default meta;
type Story = StoryObj<typeof meta>;

export const EventCard: Story = {
  args: {
    status: "rejected",
    title: "SwampHacks XI",
    description: "UF’s flagship hackathon celebrates its 11th iteration.",
    date: "Jan 25 - 26",
    location: "Newell Hall",
  },
};
