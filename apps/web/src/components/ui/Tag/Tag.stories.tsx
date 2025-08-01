/* eslint-disable */
import type { Meta } from "@storybook/react";
import { Tag, TagGroup } from ".";

const meta: Meta<typeof Example> = {
  title: "UI/Tag",
  component: TagGroup,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Example = (args: any) => (
  <TagGroup {...args} onRemove={() => {}}>
    <Tag>Computer Science</Tag>
    <Tag isDisabled>Statistics</Tag>
    <Tag>Accounting</Tag>
    <Tag>Chemistry</Tag>
  </TagGroup>
);

Example.args = {
  label: "Ice cream flavor",
  selectionMode: "single",
};
