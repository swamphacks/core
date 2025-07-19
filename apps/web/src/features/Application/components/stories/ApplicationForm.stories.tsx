/* eslint-disable */
import type { Meta } from "@storybook/react";
import { ApplicationForm } from "../ApplicationForm";

const meta: Meta<typeof ApplicationForm> = {
  title: "Event/Application",
  component: ApplicationForm,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Test = (args: any) => (
  <div className="w-screen h-screen bg-surface">
    <div className="w-full bg-surface transition-[background]">
      <ApplicationForm {...args} />
    </div>
  </div>
);

Test.args = {
  title: "SwampHacks XI Application",
  description:
    "SwampHacks is the University of Florida’s largest annual hackathon. A 36-hour tech event where students from across the country come together to build projects, learn new skills, and connect with fellow innovators.",
};
