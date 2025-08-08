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
import { useEffect, useMemo, useRef, useState, type ReactNode } from "react";
import z from "zod";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { textFieldIcons } from "@/features/FormBuilder/icons";
import { useStore } from "@tanstack/react-form";
import { questionTypeItemMap } from "./questions/createQuestionItem";

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

  const schema = questionItem.extractValidationSchemaFromItem(item);

  return !item.isRequired ? schema.optional() : schema;
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

export function build(
  formObject: FormObject,
): (props: { onSubmit?: (data: any) => void }) => ReactNode {
  const { error, data } = FormSchema.safeParse(formObject);

  if (error) {
    console.error("Invalid form json object passed to form builder.");
    throw error;
  }

  const validationSchema = getFormValidationSchema(data.content);

  return function Component({ onSubmit }) {
    const form = useAppForm({
      validators: {
        onSubmit: validationSchema,
      },
      onSubmit: ({ value }) => {
        onSubmit?.(value);
      },
    });

    const formRef = useRef<HTMLFormElement>(null);

    const formErrors = useFormErrors(form);

    // Recursively create form field components based on field type (section, layout, or question) and question types
    const buildFormContent = (content: FormItemSchemaType[]) => {
      return content.map((item) => {
        if (item.type === "question") {
          return (
            <form.AppField
              key={item.name}
              name={item.name}
              children={(field) => {
                switch (item.questionType) {
                  case QuestionTypes.shortAnswer:
                  case QuestionTypes.url:
                    return (
                      <field.TextField
                        {...item}
                        name={field.name}
                        type="text"
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
                        {...item}
                        name={field.name}
                        className="flex-1"
                        textarea
                        validationBehavior="aria"
                      />
                    );
                  case QuestionTypes.number:
                    return (
                      <field.TextField
                        {...item}
                        name={field.name}
                        type="number"
                        className="flex-1"
                        validationBehavior="aria"
                      />
                    );
                  case QuestionTypes.multipleChoice:
                    return (
                      <field.RadioField
                        {...item}
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
                    return <CheckboxField item={item} field={field} />;
                  case QuestionTypes.select:
                    return <SelectField item={item} field={field} />;
                  case QuestionTypes.multiselect:
                    return <MultiSelectField item={item} field={field} />;
                  case QuestionTypes.date:
                    return (
                      <field.DatePickerField
                        {...item}
                        name={field.name}
                        className="flex-1"
                        validationBehavior="aria"
                      />
                    );
                  case QuestionTypes.upload:
                    return (
                      <field.UploadField
                        {...item}
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
            <div className="mt-7 text-text-main" key={item.id}>
              <p className="text-xl font-medium">{item.label}</p>
              {item.description && (
                <div className="my-4">
                  <p>{item.description}</p>
                </div>
              )}
              <div className="mt-4 space-y-4">
                {buildFormContent(item.content)}
              </div>
            </div>
          );
        }

        if (item.type === "layout") {
          return (
            <FieldGroup className="gap-4" key={item.id}>
              {buildFormContent(item.content)}
            </FieldGroup>
          );
        }
      });
    };

    // Cache result so buildFormContent doesn't invoke on rerender
    const formContent = useMemo(() => buildFormContent(data.content), []);

    const errors = useStore(form.store, (state) => {
      return state.errors;
    });

    // scroll the first invalid element in the form into view
    useEffect(() => {
      if (errors.length === 0 || !formRef.current) return;

      const form = formRef.current;

      for (let i = 0; i < form.elements.length; i++) {
        const element = form.elements[i] as HTMLElement;
        const fieldName = element.getAttribute("name");

        if (fieldName) {
          if (formErrors[fieldName].length > 0) {
            const hidden =
              element.getAttribute("type") === "hidden" ||
              element.getAttribute("custom-hidden");

            if (
              (element instanceof HTMLInputElement ||
                element instanceof HTMLTextAreaElement) &&
              !hidden
            ) {
              element.focus();
            } else {
              if (hidden) {
                // if the element is hidden, use parent to scroll. this is needed for datepicker field
                element.parentElement?.scrollIntoView({ block: "center" });
                return;
              }
              element.scrollIntoView({ block: "center" });
            }

            return;
          }
        }
      }
    }, [errors]);

    return (
      <div className="w-full sm:max-w-180 mx-auto font-figtree p-2 relative">
        <div className="space-y-3 py-5 border-b-1 border-border">
          <p className="text-2xl text-text-main font-medium">
            {data.metadata.title}
          </p>
          <p className="text-text-main">{data.metadata.description}</p>
        </div>
        <Form
          className="mt-5 space-y-10"
          onSubmit={(e) => {
            e.preventDefault();
            e.stopPropagation();
            form.handleSubmit();
          }}
          validationErrors={formErrors}
          validationBehavior="aria"
          ref={formRef}
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

function useJSONData(
  item: Extract<FormItemSchemaType, { questionType: "select" | "multiselect" }>,
) {
  const [data, setData] = useState<{ id: string; name: string }[]>([]);

  useEffect(() => {
    async function fetchData() {
      if (typeof item.options === "object" && "data" in item.options) {
        const data = item.options.data;

        if (data === "schools") {
          fetch(`/public/assets/schools.json`)
            .then((res) => res.json())
            .then((schools) => {
              setData(
                schools.map((item: string) => ({ id: item, name: item })),
              );
            });
        } else if (data === "majors") {
          fetch(`/public/assets/majors.json`)
            .then((res) => res.json())
            .then((majors) => {
              setData(majors.map((item: string) => ({ id: item, name: item })));
            });
        }
      }
    }

    fetchData();
  }, []);

  return data;
}

function SelectField({
  item,
  field,
}: {
  item: Extract<FormItemSchemaType, { questionType: "select" }>;
  field: any;
}) {
  const data = useJSONData(item);

  const yearData = useMemo(() => {
    const data = [];
    if (typeof item.options === "object" && "data" in item.options) {
      if (item.options.data === "year") {
        const minYear = item.options.min!;
        const maxYear = item.options.max!;

        for (let i = minYear; i <= maxYear; i++) {
          data.push({
            id: i?.toString(),
            name: i?.toString(),
          });
        }

        return data;
      } else {
        return undefined;
      }
    } else {
      return undefined;
    }
  }, [item]);

  return item.searchable ? (
    <field.ComboBoxField
      {...item}
      name={field.name}
      className="flex-1"
      validationBehavior="aria"
      virtualized={data.length > 100}
      items={
        Array.isArray(item.options) ? item.options : yearData ? yearData : data
      }
    />
  ) : (
    <field.SelectField
      {...item}
      name={field.name}
      className="flex-1 max-h-70"
      validationBehavior="aria"
      virtualized={data.length > 100}
      items={
        Array.isArray(item.options) ? item.options : yearData ? yearData : data
      }
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
  const data = useJSONData(item);

  const transformData = (data: { id: string; name: string }[]) =>
    data.map((item) => ({ label: item.name, value: item.id }));

  return (
    <field.MultiSelectField
      {...item}
      name={field.name}
      options={Array.isArray(item.options) ? item.options : transformData(data)}
      validationBehavior="aria"
    />
  );
}

function CheckboxField({
  item,
  field,
}: {
  item: Extract<FormItemSchemaType, { questionType: "checkbox" }>;
  field: any;
}) {
  // Render a single checkbox without label. This is to handle confirmation checkboxes
  if (item.options.length === 1 && !item.label) {
    const option = item.options[0];

    return (
      <div className="flex-1">
        <Checkbox
          validationBehavior="aria"
          name={field.name}
          key={option.value}
          value={option.value}
          onChange={(val) => field.handleChange(val && [option.value])}
          className="text-base"
          isRequired={item.isRequired}
        >
          <span>
            {option.label}

            <span className="text-base text-red-500 ml-1">*</span>
          </span>
        </Checkbox>
      </div>
    );
  }

  return (
    <field.CheckboxField
      {...item}
      name={field.name}
      className="flex-1"
      validationBehavior="aria"
      isRequired={item.isRequired}
    >
      {item.options.map((option) => (
        <Checkbox key={option.value} value={option.value}>
          {option.label}
        </Checkbox>
      ))}
    </field.CheckboxField>
  );
}
