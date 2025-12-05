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
  memo,
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import z from "zod";
import { QuestionTypes } from "@/features/FormBuilder/types";
import { textFieldIcons } from "@/features/FormBuilder/icons";
import { useStore } from "@tanstack/react-form";
import { questionTypeItemMap } from "./questions/createQuestionItem";
import { parseDate } from "@internationalized/date";
import { useDebounce } from "@uidotdev/usehooks";
import TablerCornerDownRight from "~icons/tabler/corner-down-right";
import TablerCircleCheck from "~icons/tabler/circle-check";

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
export function getFormValidationSchemaAndFields(
  content: FormItemSchemaType[],
): {
  validationSchema: z.ZodType;
  fields: string[];
  fieldsMeta: Record<string, keyof typeof QuestionTypes>;
} {
  const schema: Record<string, z.ZodType<any>> = {};
  const fields: string[] = [];
  const fieldsMeta: Record<string, keyof typeof QuestionTypes> = {};

  const traverseFormContent = (content: FormItemSchemaType[]) => {
    // For every question, attempt to parse its validation schema from the `validation` field if exist.
    for (const item of content) {
      if (item.type === "question") {
        schema[item.name] = getFormItemValidationSchema(item);
        fields.push(item.name);
        fieldsMeta[item.name] = item.questionType;

        if (item.questionType === QuestionTypes.select) {
          if (item.hasOther) {
            fields.push(`${item.name}-other`);
            fieldsMeta[`${item.name}-other`] = QuestionTypes.shortAnswer;
            schema[`${item.name}-other`] = getFormItemValidationSchema({
              ...item,
              isRequired: false,
            });
          }
        }
      } else {
        traverseFormContent(item.content);
      }
    }
  };

  traverseFormContent(content);

  return {
    validationSchema: z.object(schema),
    fields,
    fieldsMeta,
  };
}

// this function ensure that default values are in the proper shape according to their question's zod schema
function transformDefaultValues(
  defaultValues: Record<string, any>,
  fieldsMeta: Record<string, keyof typeof QuestionTypes>,
): Record<string, any> {
  for (const field in defaultValues) {
    if (fieldsMeta[field] === QuestionTypes.checkbox) {
      if (typeof defaultValues[field] === "string") {
        defaultValues[field] = defaultValues[field].split(",");
      }
    } else if (fieldsMeta[field] === QuestionTypes.multiselect) {
      if (typeof defaultValues[field] === "string") {
        defaultValues[field] = defaultValues[field].split(",");
      }
    } else if (fieldsMeta[field] === QuestionTypes.number) {
      defaultValues[field] = defaultValues[field].toString();
    }
  }

  return defaultValues;
}

export interface FormProps {
  onSubmit?: (data: Record<string, any>) => Promise<void>;
  onNewAttachments?: (data: Record<string, File[]>) => void;
  onChangeDelayMs?: number;
  onChange?: (formValues: Record<string, any>) => void;
  defaultValues?: Record<string, any>;
  isInvalid?: boolean;
  isSubmitted?: boolean;
  isSubmitting?: boolean;
  SubmitSuccessComponent?: React.ComponentType;
  renderFormHeader?: (metadata: {
    title: string;
    description?: string | undefined;
  }) => React.ReactNode;
}

export function build(formObject: FormObject): {
  Form: (props: FormProps) => ReactNode;
  fields: string[];
  fieldsMeta: Record<string, keyof typeof QuestionTypes>;
  defaultFieldValues: Record<string, undefined>;
} {
  const { error, data } = FormSchema.safeParse(formObject);

  if (error) {
    console.error("Invalid form json object passed to form builder.");
    throw error;
  }

  const { validationSchema, fields, fieldsMeta } =
    getFormValidationSchemaAndFields(data.content);

  const defaultFieldValues: Record<string, undefined> = {};

  for (const key of fields) {
    defaultFieldValues[key] = undefined;
  }

  return {
    fields,
    fieldsMeta,
    defaultFieldValues,
    Form: memo(function Component({
      onSubmit,
      onNewAttachments,
      onChange,
      defaultValues = {},
      onChangeDelayMs = 5000,
      isInvalid = false,
      isSubmitted: isSubmittedProp = false,
      isSubmitting = false,
      SubmitSuccessComponent = FallbackSubmitSuccessComponent,
      renderFormHeader,
    }) {
      const [otherFields, setOtherFields] = useState<string[]>([]);

      const transformedDefaultValues = useMemo(
        () => transformDefaultValues(defaultValues, fieldsMeta),
        [defaultValues],
      );

      const form = useAppForm({
        defaultValues: transformedDefaultValues,
        validators: {
          // @ts-expect-error - validationSchema type doesn't match expected validator type
          onSubmit: validationSchema,
        },
        onSubmit: async ({ value }) => {
          await onSubmit?.(value);
        },
      });

      const isDirty = useDebounce(form.state.isDirty, onChangeDelayMs);
      const formValues = useDebounce(
        form.state.values,
        onChangeDelayMs,
      ) as Record<string, any>;

      const previousValuesRef = useRef<string>("");

      useEffect(() => {
        if (isSubmitting || isSubmittedProp) return;

        if (isDirty && Object.keys(formValues).length >= 1) {
          const currentValues = JSON.stringify({
            ...defaultFieldValues,
            ...formValues,
          });

          if (previousValuesRef.current !== currentValues) {
            previousValuesRef.current = currentValues;
            onChange?.({
              ...defaultFieldValues,
              ...formValues,
            });
          }
        }
      }, [isDirty, formValues, isSubmitting, isSubmittedProp, onChange]);

      useEffect(() => {
        const newOtherFields = [];
        for (const field in form.state.values) {
          if (
            form.state.values[field] === "other" &&
            fieldsMeta[field] === QuestionTypes.select
          ) {
            newOtherFields.push(field);
          }
        }
        // Only update if arrays differ
        if (
          newOtherFields.length !== otherFields.length ||
          !newOtherFields.every((field) => otherFields.includes(field))
        ) {
          setOtherFields(newOtherFields);
        }
      }, [form.state.isDirty, form.state.values, otherFields, fieldsMeta]);

      const formRef = useRef<HTMLFormElement>(null);

      const formErrors = useFormErrors(form);

      // Recursively create form field components based on field type (section, layout, or question) and question types
      const buildFormContent = useCallback(
        (content: FormItemSchemaType[]) => {
          return content.map((item) => {
            if (item.type === "question") {
              if (item.questionType === QuestionTypes.select) {
                if (otherFields.includes(item.name)) {
                  return (
                    <div
                      key={`${item.name}-1`}
                      className="flex-1 flex flex-col gap-1"
                    >
                      <form.AppField
                        key={item.name}
                        name={item.name}
                        children={(field) => {
                          return <SelectField item={item} field={field} />;
                        }}
                      />
                      <div className="flex items-center gap-2">
                        <TablerCornerDownRight className="opacity-20 text-xl mb-2" />
                        <form.AppField
                          key={`${item.name}-other`}
                          name={`${item.name}-other`}
                          children={(field) => {
                            return (
                              <field.TextField
                                placeholder="Enter value"
                                aria-label={`${item.name}-other`}
                                name={field.name}
                                defaultValue={field.state.value}
                                type="text"
                                className="flex-1"
                                validationBehavior="aria"
                              />
                            );
                          }}
                        />
                      </div>
                    </div>
                  );
                }

                return (
                  <form.AppField
                    key={item.name}
                    name={item.name}
                    children={(field) => {
                      return <SelectField item={item} field={field} />;
                    }}
                  />
                );
              }
              return (
                <form.AppField
                  key={item.name}
                  name={item.name}
                  children={(field) => {
                    const fieldValue = field.state.value as string | undefined;

                    switch (item.questionType) {
                      case QuestionTypes.shortAnswer:
                      case QuestionTypes.url:
                        return (
                          <field.TextField
                            {...item}
                            name={field.name}
                            defaultValue={fieldValue}
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
                            defaultValue={fieldValue}
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
                            defaultValue={fieldValue}
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
                            defaultValue={fieldValue}
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
                      case QuestionTypes.multiselect:
                        return <MultiSelectField item={item} field={field} />;
                      case QuestionTypes.date:
                        return (
                          <field.DatePickerField
                            {...item}
                            name={field.name}
                            className="flex-1"
                            validationBehavior="aria"
                            defaultValue={parseDate(fieldValue ?? "")} // TODO: verify correctness
                          />
                        );
                      case QuestionTypes.upload:
                        return (
                          <field.UploadField
                            {...item}
                            name={field.name}
                            validationBehavior="aria"
                            defaultValue={field.state.value as any}
                            onChange={(files) => {
                              field.handleChange(
                                files.length === 0 ? null : files,
                              );
                            }}
                            onNewFiles={(newFiles) => {
                              onNewAttachments?.({
                                [field.name]: newFiles,
                              });
                            }}
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
        },
        [otherFields, form, onNewAttachments],
      );

      // Cache result so buildFormContent doesn't invoke on rerender
      const formContent = useMemo(
        () => buildFormContent(data.content),
        [buildFormContent],
      );

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
      }, [errors, formErrors]);

      const isSubmittedState = useStore(form.store, (state) => {
        return state.isSubmitted;
      });

      const isSubmitted = (isSubmittedState || isSubmittedProp) && !isInvalid;

      // TODO: find a better way to handle this, maybe just let the parent component handle the form component instead of doing it in here
      const formContainer = useMemo(() => {
        return (
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
        );
      }, [formErrors, formContent, form]);

      const formHeader = useMemo(() => {
        if (renderFormHeader) {
          return renderFormHeader(data.metadata);
        } else {
          return (
            <div className="space-y-3 py-5 border-b-1 border-border">
              <p className="text-2xl text-text-main font-medium">
                {data.metadata.title}
              </p>
              <p className="text-text-main">{data.metadata.description}</p>
            </div>
          );
        }
      }, [renderFormHeader]);

      return (
        <div className="w-full sm:max-w-180 mx-auto font-figtree p-2 relative">
          {formHeader}
          {isSubmitted ? (
            <div className="mt-3">
              <SubmitSuccessComponent />
            </div>
          ) : (
            formContainer
          )}
        </div>
      );
    }),
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

        switch (data) {
          case "schools": {
            const schoolsRes = await fetch(`/assets/schools.json`);
            const schools = await schoolsRes.json();
            setData(schools.map((item: string) => ({ id: item, name: item })));
            break;
          }
          case "majors": {
            const majorsRes = await fetch(`/assets/majors.json`);
            const majors = await majorsRes.json();
            setData(majors.map((item: string) => ({ id: item, name: item })));
            break;
          }
          case "countries": {
            const countriesRes = await fetch(`/assets/countries.json`);
            const countries = await countriesRes.json();
            setData(
              countries.map(
                ({ name, code }: { name: string; code: string }) => ({
                  id: code,
                  name,
                }),
              ),
            );
            break;
          }
          default:
            break;
        }
      }
    }

    fetchData();
  }, [item.options]);

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
      defaultSelectedKey={field.state.value}
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
      defaultSelectedKey={field.state.value}
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

  let defaultValue = field.state.value;

  const mapFn = (item: any) => {
    if (typeof item === "object" && item["label"] && item["value"]) {
      return item;
    }

    return {
      label: item,
      value: item,
    };
  };

  if (typeof defaultValue === "string") {
    defaultValue = defaultValue.split(",").map(mapFn);
  } else {
    defaultValue = Array.isArray(defaultValue) ? defaultValue.map(mapFn) : [];
  }

  return (
    <field.MultiSelectField
      {...item}
      name={field.name}
      options={Array.isArray(item.options) ? item.options : transformData(data)}
      validationBehavior="aria"
      defaultValue={defaultValue}
      onChange={(data: any) =>
        field.handleChange(data.map((item: any) => item.value))
      }
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
          onChange={(val) => field.setValue(val ? true : null)}
          className="text-base"
          isRequired={item.isRequired}
          defaultSelected={field.state.value}
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
      defaultValue={field.state.value}
    >
      {item.options.map((option) => (
        <Checkbox key={option.value} value={option.value}>
          {option.label}
        </Checkbox>
      ))}
    </field.CheckboxField>
  );
}

function FallbackSubmitSuccessComponent() {
  return (
    <div className="flex items-center gap-2 py-3 pl-3 bg-badge-bg-accepted/50 text-badge-text-accepted font-medium">
      <TablerCircleCheck />
      <p>Thank you! Your application has been received.</p>
    </div>
  );
}
