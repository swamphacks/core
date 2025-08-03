import type { Meta } from "@storybook/react";
import { Form, useAppForm } from ".";

const meta: Meta<typeof Form> = {
  title: "UI/Form",
  component: Form,
  parameters: {
    layout: "centered",
  },
  tags: ["autodocs"],
};

export default meta;

export const Example = () => {
  const form = useAppForm({
    onSubmit: ({ value }) => {
      console.log(value);
    },
  });

  return (
    <Form
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
    >
      <form.AppField
        name="email"
        children={(field) => (
          <field.TextField
            label="Email"
            name={field.name}
            type="email"
            autoComplete="off"
            className="flex-1"
            placeholder="Enter your email"
            isRequired={true}
          />
        )}
      />

      <form.AppField
        name="date"
        children={(field) => (
          <field.DatePickerField
            label="Date"
            name={field.name}
            className="flex-1"
            isRequired={true}
          />
        )}
      />

      <div className="flex gap-2">
        <form.AppForm>
          <form.SubmitButton label="Submit" />
        </form.AppForm>
      </div>
    </Form>
  );
};
