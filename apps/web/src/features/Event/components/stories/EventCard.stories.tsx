import type { Meta, StoryObj } from "@storybook/react";
import { EventCard as EventCardComp } from "../EventCard";
import applicationStatus from "../../applicationStatus";

const meta = {
  component: EventCardComp,
  title: "Event/Event Card",
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
    eventId: "test id",
    status: "rejected",
    title: "SwampHacks XI",
    description: "UF’s flagship hackathon celebrates its 11th iteration.",
    date: "Jan 25 - 26",
    location: "Newell Hall",
    application_close: new Date("2024-01-15T23:59:59"),
    application_open: new Date("2023-11-01T00:00:00"),
    end_time: new Date("2024-01-26T20:00:00"),
    start_time: new Date("2024-01-25T09:00:00"),
    banner:
      "https://media.istockphoto.com/id/974238866/photo/audience-listens-to-the-lecturer-at-the-conference.jpg?s=612x612&w=0&k=20&c=p_BQCJWRQQtZYnQlOtZMzTjeB_csic8OofTCAKLwT0M=",
  },
};
