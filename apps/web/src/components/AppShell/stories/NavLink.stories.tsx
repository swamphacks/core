/* This is bugged out brother */
import type { Meta, StoryObj } from "@storybook/react";
import TablerLayoutCollage from "~icons/tabler/layout-collage";
import { NavLink } from "../NavLink";
import {
  createRootRouteWithContext,
  createRouter,
  RouterProvider,
} from "@tanstack/react-router";
import type { ReactNode } from "react";

// Define a dummy root component that wraps Story
function withTanstackRouter(StoryComponent: () => ReactNode) {
  const RootRoute = createRootRouteWithContext()({
    component: StoryComponent,
  });

  const router = createRouter({
    routeTree: RootRoute,
  });

  return <RouterProvider router={router} />;
}

const meta = {
  component: NavLink,
  title: "AppShell/NavLink",
  tags: ["autodocs"],
  decorators: [(Story) => withTanstackRouter(Story)],
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
    leftSection: (
      <TablerLayoutCollage strokeWidth={0} className="w-4 aspect-square" />
    ),
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

export const Expandable: Story = {
  args: {
    label: "Mother",
    children: (
      <>
        <NavLink label="Child 1" />
        <NavLink label="Child 2" />
      </>
    ),
  },
};

export const WithDescription: Story = {
  args: {
    label: "Resume Review",
    description: "67/67 (100%) resumes reviewed.",
  },
};
