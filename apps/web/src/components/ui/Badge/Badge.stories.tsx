import type { Meta, StoryObj } from "@storybook/react";
import { Badge } from ".";

const meta = {
  component: Badge,
  title: "UI/Badge",
  tags: ["autodocs"],
} satisfies Meta<typeof Badge>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    children: "Badge",
    type: "default",
    size: "sm",
  },
};

export const DefaultWithIcon: Story = {
  args: {
    children: (
      <>
        <svg
          className="w-2.5 h-2.5"
          viewBox="0 0 11 11"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M2.61248 2.61248L8.38748 8.38748M1.375 5.5C1.375 6.0417 1.4817 6.5781 1.689 7.07857C1.8963 7.57904 2.20014 8.03377 2.58318 8.41682C2.96623 8.79986 3.42096 9.1037 3.92143 9.311C4.4219 9.5183 4.9583 9.625 5.5 9.625C6.0417 9.625 6.5781 9.5183 7.07857 9.311C7.57904 9.1037 8.03377 8.79986 8.41682 8.41682C8.79986 8.03377 9.1037 7.57904 9.311 7.07857C9.5183 6.5781 9.625 6.0417 9.625 5.5C9.625 4.9583 9.5183 4.4219 9.311 3.92143C9.1037 3.42096 8.79986 2.96623 8.41682 2.58318C8.03377 2.20014 7.57904 1.8963 7.07857 1.689C6.5781 1.4817 6.0417 1.375 5.5 1.375C4.9583 1.375 4.4219 1.4817 3.92143 1.689C3.42096 1.8963 2.96623 2.20014 2.58318 2.58318C2.20014 2.96623 1.8963 3.42096 1.689 3.92143C1.4817 4.4219 1.375 4.9583 1.375 5.5Z"
            stroke="#E7000B"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
        Badge
      </>
    ),
    size: "sm",
  },
};

export const DefaultMediumSize: Story = {
  args: {
    children: <>Badge</>,
    size: "md",
  },
};
