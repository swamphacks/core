import type { Meta, StoryObj } from "@storybook/react";
import { Avatar } from "./Avatar";

const meta = {
  component: Avatar,
  title: "UI/Avatar",
  tags: ["autodocs"],
} satisfies Meta<typeof Avatar>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    src: "https://preview.redd.it/miku-is-in-your-wifi-now-shes-in-your-screen-and-licking-it-v0-e3jtii4ohtmf1.png?width=256&auto=webp&s=6206333929d76a8f4e5fda19663f9f7634dfaebe",
  },
};

export const DefaultWithStatus: Story = {
  args: {
    src: "https://preview.redd.it/miku-is-in-your-wifi-now-shes-in-your-screen-and-licking-it-v0-e3jtii4ohtmf1.png?width=256&auto=webp&s=6206333929d76a8f4e5fda19663f9f7634dfaebe",
    status: "away",
  },
};

export const ExtraLargeWithStatus: Story = {
  args: {
    src: "https://preview.redd.it/miku-is-in-your-wifi-now-shes-in-your-screen-and-licking-it-v0-e3jtii4ohtmf1.png?width=256&auto=webp&s=6206333929d76a8f4e5fda19663f9f7634dfaebe",
    status: "online",
    size: "xl",
  },
};

export const Fallback: Story = {
  args: {
    fallback: "Swamp Hacks",
  },
};

export const Empty: Story = {
  args: {},
};

export const Square: Story = {
    args: {
        shape: "square"
    },
  };
  
