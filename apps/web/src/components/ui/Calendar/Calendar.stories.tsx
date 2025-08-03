import { Calendar } from ".";

import type { Meta } from "@storybook/react";

const meta: Meta<typeof Calendar> = {
  title: "UI/Calendar",
  component: Calendar,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

/* eslint-disable @typescript-eslint/no-explicit-any */
export const Example = (args: any) => (
  <Calendar aria-label="Event date" {...args} />
);
