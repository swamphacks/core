import type { Meta, StoryObj } from "@storybook/react";
import { AvatarStack } from "./AvatarStack";

const meta = {
  component: AvatarStack,
  title: "UI/AvatarStack",
  tags: ["autodocs"],
} satisfies Meta<typeof AvatarStack>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
    args: {
        avatars: [
            { src: "https://i.pinimg.com/1200x/b3/33/b9/b333b9da782cb6874da5002f185a80ce.jpg" },
            { src: "https://i.pinimg.com/736x/d9/fa/c3/d9fac3d552659e888bbecbd20b12f160.jpg" },
            { fallback: "Alexander Wang" },
            { fallback: "Phoenix" },
            { fallback: "Hieu" },
            { fallback: "Hugo" },
            { fallback: "Stanley" },
            { fallback: "Kayla" }
        ],
        max: 3,
        overlap: "medium"
    }
}
