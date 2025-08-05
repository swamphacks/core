import { Form, useAppForm, useFormErrors } from "@/components/Form";
import { Checkbox } from "@/components/ui/Checkbox";
import { FieldGroup } from "@/components/ui/Field";
import { Radio } from "@/components/ui/Radio";
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
import { useEffect, useMemo, useState, type ReactNode } from "react";
import z from "zod";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { textFieldIcons } from "@/features/FormBuilder/icons";
// import { FormValidationContext } from "react-stately";
// import { FormContext } from "react-aria-components";
// import { useStore } from "@tanstack/react-form";

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

    // const formRef = useRef<HTMLFormElement>(null);

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
                        validationBehavior="aria"
                        icon={
                          item.iconName
                            ? textFieldIcons[item.iconName]
                            : undefined
                        }
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
                        validationBehavior="aria"
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
                        validationBehavior="aria"
                      />
                    );
                  case QuestionTypes.multipleChoice:
                    return (
                      <field.RadioField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                        validationBehavior="aria"
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
                        validationBehavior="aria"
                        isRequired={true}
                      >
                        {item.options.map((option) => (
                          <Checkbox key={option.value} value={option.value}>
                            {option.label}
                          </Checkbox>
                        ))}
                      </field.CheckboxField>
                    );
                  case QuestionTypes.select:
                    return (
                      <SelectField
                        item={
                          newItem as Extract<
                            FormItemSchemaType,
                            { questionType: "select" }
                          >
                        }
                        field={field}
                      />
                    );
                  case QuestionTypes.multiselect:
                    return (
                      <MultiSelectField
                        item={
                          newItem as Extract<
                            FormItemSchemaType,
                            { questionType: "multiselect" }
                          >
                        }
                        field={field}
                      />
                    );
                  case QuestionTypes.date:
                    return (
                      <field.DatePickerField
                        {...newItem}
                        name={field.name}
                        className="flex-1"
                        validationBehavior="aria"
                      />
                    );
                  case QuestionTypes.upload:
                    return (
                      <field.UploadField
                        {...newItem}
                        name={field.name}
                        validationBehavior="aria"
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

    // const errors = useStore(form.store, (state) => {
    //   return state.errors;
    // });

    // TODO: scroll to the element that has an error. what I have below only works for certain components, doesnt work for Select or Radio fields
    // useEffect(() => {
    //   if (errors.length === 0) return;

    //   if (!formRef.current) return;

    //   const form = formRef.current;

    //   for (let i = 0; i < form.elements.length; i++) {
    //     let element = form.elements[i];
    //     // console.log(element);
    //     if (
    //       element.getAttribute("data-invalid") ||
    //       element.getAttribute("data-custom-invalid")
    //     ) {
    //       element.focus();
    //       // console.log(element);
    //       return;
    //     }
    //   }
    // }, [errors]);

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
          validationBehavior="aria"
        >
          <div className="space-y-6">{formContent}</div>
          <form.AppForm>
            <form.SubmitButton label="Submit" />
          </form.AppForm>
        </Form>
        {/* <form
          className="mt-5"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          ref={formRef}
        >
          <FormContext.Provider value={{ validationBehavior: "aria" }}>
            <FormValidationContext.Provider value={formErrors}>
              <div className="space-y-6">{formContent}</div>
              <form.AppForm>
                <form.SubmitButton label="Submit" />
              </form.AppForm>
            </FormValidationContext.Provider>
          </FormContext.Provider>
        </form> */}
      </div>
    );
  };
}

function SelectField({
  item,
  field,
}: {
  item: Extract<FormItemSchemaType, { questionType: "select" }>;
  field: any;
}) {
  const [data, setData] = useState<{ id: string; name: string }[]>([]);

  useEffect(() => {
    if (typeof item.options === "object" && "data" in item.options) {
      const data = item.options.data;

      // TODO: where should we put these json files? is it already cached?
      if (data === "schools") {
        import(`./scripts/schools.json`).then((schools) => {
          setData(schools.default);
        });
      } else if (data === "majors") {
        import(`./scripts/majors.json`).then((majors) => {
          setData(majors.default);
        });
      }
    }
  }, []);

  return item.searchable ? (
    <field.ComboBoxField
      {...item}
      name={field.name}
      className="flex-1"
      validationBehavior="aria"
      virtualized={data.length > 100}
      items={Array.isArray(item.options) ? item.options : data}
    />
  ) : (
    <field.SelectField
      {...item}
      name={field.name}
      className="flex-1"
      validationBehavior="aria"
      virtualized={data.length > 100}
      items={Array.isArray(item.options) ? item.options : data}
    />
  );
}

function MultiSelectField({
  item,
  field,
}: {
  item: Extract<FormItemSchemaType, { questionType: "multiselect" }>;
  field: any;
}) {
  const [data, setData] = useState<{ label: string; value: string }[]>([]);

  const transformData = (data: { id: string; name: string }[]) =>
    data.map((item) => ({ label: item.name, value: item.id }));

  useEffect(() => {
    if (typeof item.options === "object" && "data" in item.options) {
      const data = item.options.data;

      if (data === "schools") {
        import(`./scripts/schools.json`).then((schools) => {
          setData(transformData(schools.default));
        });
      } else if (data === "majors") {
        import(`./scripts/majors.json`).then((majors) => {
          setData(transformData(majors.default));
        });
      }
    }
  }, []);

  return (
    <field.MultiSelectField
      {...item}
      name={field.name}
      options={Array.isArray(item.options) ? item.options : data}
      validationBehavior="aria"
    />
  );
}
