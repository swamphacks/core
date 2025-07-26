import { Form, useAppForm, useFormErrors } from "@/components/Form";
import { Checkbox } from "@/components/ui/Checkbox";
import { FieldGroup } from "@/components/ui/Field";
import { Radio } from "@/components/ui/Radio";
import { SelectItem } from "@/components/ui/Select";
import {
  type FormItemSchemaType,
  type FormObject,
  type FormQuestionItemSchemaType,
  FormSchema,
} from "@/features/FormBuilder/formSchema";
import {
  CheckboxQuestion,
  createQuestionItem,
  DateQuestion,
  MultipleChoiceQuestion,
  MultiSelectQuestion,
  NumberQuestion,
  ParagraphQuestion,
  SelectQuestion,
  ShortAnswerQuestion,
  UploadQuestion,
} from "@/features/FormBuilder/questions";
import { useMemo, type ReactNode } from "react";
import z from "zod";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { ComboBoxItem } from "@/components/ui/ComboBox";

const questionTypeItemMap = {
  [QuestionTypes.shortAnswer]: ShortAnswerQuestion,
  [QuestionTypes.paragraph]: ParagraphQuestion,
  [QuestionTypes.number]: NumberQuestion,
  [QuestionTypes.multipleChoice]: MultipleChoiceQuestion,
  [QuestionTypes.checkbox]: CheckboxQuestion,
  [QuestionTypes.select]: SelectQuestion,
  [QuestionTypes.multiselect]: MultiSelectQuestion,
  [QuestionTypes.upload]: UploadQuestion,
  [QuestionTypes.date]: DateQuestion,
  // Use `satisfies` here so `questionTypeItemMap` doesn't actually take on this Record type, which can cause type errors in `getFormItemValidationSchema`
} satisfies Record<
  keyof typeof QuestionTypes,
  ReturnType<typeof createQuestionItem>
>;

export function getFormItemValidationSchema(
  item: FormQuestionItemSchemaType,
): z.ZodType {
  const questionItem = questionTypeItemMap[item.questionType];

  if (!questionItem) {
    /* eslint-disable @typescript-eslint/no-explicit-any */
    throw new Error(
      `Form builder error: unknown question type or missing validation parser for question type: ${(item as any).questionType}`,
    );
  }

  /* eslint-disable @typescript-eslint/no-explicit-any */
  const schema = questionItem.extractValidationSchemaFromItem(item as any); // Not sure how to handle the type for `item` here so any it is.

  return !item.required ? schema.optional() : schema;
}

// https://stackoverflow.com/a/78818575
// Dynamically build a zod schema from the validation field
export function getFormValidationSchema(
  content: FormItemSchemaType[],
): z.ZodType {
  const schema: Record<string, z.ZodType<any>> = {};

  const traverseFormContent = (content: FormItemSchemaType[]) => {
    // For every question, attempt to parse its validation schema from the `validation` field if exist.
    for (const item of content) {
      if (item.type === "question") {
        schema[item.name] = getFormItemValidationSchema(item);
      } else {
        traverseFormContent(item.content);
      }
    }
  };

  traverseFormContent(content);

  return z.object(schema);
}

// TODO:
// Better error handling
// json file upload?
// file upload (resumes)
// handle input icons?

// add date picker
// add file upload

export function build(formObject: FormObject): () => ReactNode {
  const { error, data } = FormSchema.safeParse(formObject);

  if (error) {
    console.error("Invalid form json object passed to form builder.");
    throw error;
  }

  const validationSchema = getFormValidationSchema(data.content);

  return function Component() {
    const form = useAppForm({
      validators: {
        onSubmit: validationSchema,
      },
      onSubmit: ({ value }) => {
        // TODO: pass in custom on submit prop
        console.log(value);
      },
    });

    const formErrors = useFormErrors(form);

    // Recursively create form field components based on field type (section, layout, or question) and question types
    const buildFormContent = (content: FormItemSchemaType[]) => {
      return content.map((item) => {
        if (item.type === "question") {
          return (
            <form.AppField
              key={item.id}
              name={item.name}
              children={(field) => {
                const newItem = {
                  ...item,
                  isRequired: item.required, // rename `required` to `isRequired`
                };

                // @ts-ignore
                delete newItem["required"];

                switch (item.questionType) {
                  case QuestionTypes.shortAnswer:
                    return (
                      <field.TextField
                        {...newItem}
                        name={field.name}
                        type="text"
                        autoComplete="off"
                        className="flex-1"
                      />
                    );
                  case QuestionTypes.paragraph:
                    return (
                      <field.TextField
                        {...newItem}
                        name={field.name}
                        autoComplete="off"
                        className="flex-1"
                        textarea
                      />
                    );
                  case QuestionTypes.number:
                    return (
                      <field.TextField
                        {...newItem}
                        name={field.name}
                        type="number"
                        autoComplete="off"
                        className="flex-1"
                      />
                    );
                  case QuestionTypes.multipleChoice:
                    return (
                      <field.RadioField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                      >
                        {item.options.map((option) => (
                          <Radio key={option.value} value={option.value}>
                            {option.label}
                          </Radio>
                        ))}
                      </field.RadioField>
                    );
                  case QuestionTypes.checkbox:
                    return (
                      <field.CheckboxField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                      >
                        {item.options.map((option) => (
                          <Checkbox key={option.value} value={option.value}>
                            {option.label}
                          </Checkbox>
                        ))}
                      </field.CheckboxField>
                    );
                  case QuestionTypes.select:
                    return item.searchable ? (
                      <field.ComboBoxField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                      >
                        {item.options.map((option) => (
                          <ComboBoxItem key={option.value} id={option.value}>
                            {option.label}
                          </ComboBoxItem>
                        ))}
                      </field.ComboBoxField>
                    ) : (
                      <field.SelectField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                      >
                        {item.options.map((option) => (
                          <SelectItem key={option.value} id={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </field.SelectField>
                    );
                  case QuestionTypes.multiselect:
                    return (
                      <field.MultiSelectField
                        {...newItem}
                        name={field.name}
                        options={item.options}
                      />
                    );
                  case QuestionTypes.date:
                    return (
                      <field.DatePickerField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                      />
                    );
                  case QuestionTypes.upload:
                    return (
                      <field.UploadField
                        {...newItem}
                        name={field.name}
                        multiple
                        // className="flex-1"
                      />
                    );
                  default:
                    break;
                }
              }}
            />
          );
        }

        if (item.type === "section") {
          return (
            <div key={item.label}>
              <p className="text-xl font-medium">{item.label}</p>
              <div className="mt-3 space-y-4">
                {buildFormContent(item.content)}
              </div>
            </div>
          );
        }

        if (item.type === "layout") {
          return (
            // TODO: remove math.random()
            <FieldGroup className="gap-4" key={Math.random()}>
              {buildFormContent(item.content)}
            </FieldGroup>
          );
        }
      });
    };

    // Cache result so buildFormContent doesn't invoke on rerender
    const formContent = useMemo(() => buildFormContent(data.content), []);

    return (
      <div className="w-full sm:max-w-180 mx-auto font-figtree pb-2 p-2">
        <div className="space-y-3 py-5 border-b-1 border-border">
          <p className="text-2xl text-text-main">{data.metadata.title}</p>
          <p className="text-text-secondary">{data.metadata.description}</p>
        </div>
        <Form
          className="mt-5"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          validationErrors={formErrors}
        >
          <div className="space-y-6">{formContent}</div>
          <form.AppForm>
            <form.SubmitButton label="Submit" />
          </form.AppForm>
        </Form>
      </div>
    );
  };
}
