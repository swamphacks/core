import type { Meta } from "@storybook/react";
import { build } from "../build";
import data from "./example.json";

const meta: Meta = {
  title: "Form Builder/Example",
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Test = () => {
  const Form = build(data);
  return <Form />;
};
