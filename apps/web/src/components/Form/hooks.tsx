// https://tanstack.com/form/latest/docs/framework/react/guides/form-composition

import {
  createFormHook,
  FormApi,
  useStore,
  type StandardSchemaV1Issue,
} from "@tanstack/react-form";
import { lazy } from "react";
import { fieldContext, formContext } from "./formContext.ts";
import { SubmitButton } from "@/components/Form/fields/SubmitButton";

const TextField = lazy(() => import("./fields/TextField"));
const CheckboxField = lazy(() => import("./fields/CheckboxField"));
const ComboBoxField = lazy(() => import("./fields/ComboBoxField"));
const MultiSelectField = lazy(() => import("./fields/MultiSelectField"));
const RadioField = lazy(() => import("./fields/RadioField"));
const SelectField = lazy(() => import("./fields/SelectField"));
const DatePickerField = lazy(() => import("./fields/DatePickerField.tsx"));

export const { useAppForm, withForm } = createFormHook({
  fieldComponents: {
    TextField,
    CheckboxField,
    ComboBoxField,
    MultiSelectField,
    RadioField,
    SelectField,
    DatePickerField,
  },
  formComponents: {
    SubmitButton,
  },
  fieldContext,
  formContext,
});

/* eslint-disable @typescript-eslint/no-explicit-any */
export function useFormErrors<T>(
  form: FormApi<T, any, any, any, any, any, any, any, any, any>, // I can't find anything online or on the docs on how to handle this
): Record<string, string[]> {
  const fieldMeta = useStore(form.store, (state) => {
    return state.fieldMeta;
  });
  const formErrors: Record<string, string[]> = {};

  const fields = Object.keys(fieldMeta);

  fields.forEach((field) => {
    formErrors[field] = fieldMeta[field as keyof typeof fieldMeta].errors.map(
      (error: StandardSchemaV1Issue) => error.message,
    );
  });

  return formErrors;
}
