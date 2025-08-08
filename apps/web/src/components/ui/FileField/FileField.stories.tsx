import type { Meta } from "@storybook/react";
import { FileField } from ".";

const meta: Meta<typeof FileField> = {
  title: "UI/FileField",
  component: FileField,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Example = () => {
  return <FileField label="Upload" />;
};
