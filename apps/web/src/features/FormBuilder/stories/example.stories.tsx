import type { Meta } from "@storybook/react";
import { build } from "../build";
import data from "./example.json";
import applicationFormData from "./applicationFormExample.json";

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

export const ApplicationForm = () => {
  const Form = build(applicationFormData);
  return (
    <div className="w-full h-screen bg-surface">
      <div className="w-full bg-surface transition-[background] sm:pl-[calc(100vw-100%)]">
        <Form onSubmit={(data) => console.log(data)} />
      </div>
    </div>
  );
};
ApplicationForm.parameters = {
  layout: "fullscreen",
};
