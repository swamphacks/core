import type { Meta, StoryObj } from "@storybook/react";
import { NavLink } from "./NavLink";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import TablerChevronRight from "~icons/tabler/chevron-right";

const meta = {
  component: NavLink,
  title: "AppShell/NavLink",
  tags: ["autodocs"],
  argTypes: {
    href: {
      control: { type: "text" },
      description: "The URL the link points to",
      defaultValue: "https://core.apidocumentation.com",
    },
    label: {
      control: { type: "text" },
      description: "The text displayed in the link",
      defaultValue: "Dashboard",
    },
    active: {
      control: { type: "boolean" },
      description: "Indicates if the link is active",
      defaultValue: false,
    },
  },
  args: {
    href: "https://core.apidocumentation.com",
    label: "Dashboard",
    active: true,
    leftSection: <TablerLayoutCollage strokeWidth={0} className="w-5 h-5" />,
    rightSection: <TablerChevronRight className="w-5 h-5" />,
  },
} satisfies Meta<typeof NavLink>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    label: "Dashboard",
    href: "https://core.apidocumentation.com",
  },
};
