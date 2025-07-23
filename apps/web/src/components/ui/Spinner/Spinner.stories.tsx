import type { Meta } from "@storybook/react";
import { Spinner } from ".";

const meta: Meta<typeof Spinner> = {
  title: "UI/Spinner",
  component: Spinner,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Example = () => <Spinner />;
