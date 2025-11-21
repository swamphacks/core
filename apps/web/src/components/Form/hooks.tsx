// https://tanstack.com/form/latest/docs/framework/react/guides/form-composition

import {
  createFormHook,
  FormApi,
  useStore,
  type StandardSchemaV1Issue,
} from "@tanstack/react-form";
import { lazy, useRef } from "react";
import { fieldContext, formContext } from "./formContext.ts";
import { SubmitButton } from "@/components/Form/fields/SubmitButton";

const TextField = lazy(() => import("./fields/TextField"));
const CheckboxField = lazy(() => import("./fields/CheckboxField"));
const ComboBoxField = lazy(() => import("./fields/ComboBoxField"));
const MultiSelectField = lazy(() => import("./fields/MultiSelectField"));
const RadioField = lazy(() => import("./fields/RadioField"));
const SelectField = lazy(() => import("./fields/SelectField"));
const DatePickerField = lazy(() => import("./fields/DatePickerField"));
const UploadField = lazy(() => import("./fields/UploadField"));

export const { useAppForm, withForm } = createFormHook({
  fieldComponents: {
    TextField,
    CheckboxField,
    ComboBoxField,
    MultiSelectField,
    RadioField,
    SelectField,
    DatePickerField,
    UploadField,
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
  // const formErrors: Record<string, string[]> = {};
  const prevRef = useRef<Record<string, string[]>>({});

  // const fields = Object.keys(fieldMeta);

  const newErrors: Record<string, string[]> = {};
  for (const field of Object.keys(fieldMeta)) {
    newErrors[field] = fieldMeta[field as keyof typeof fieldMeta].errors.map(
      (error: StandardSchemaV1Issue) => error.message,
    );
  }

  const prev = prevRef.current;

  // shallow compare keys and arrays of messages
  let changed = false;
  const prevKeys = Object.keys(prev);
  const newKeys = Object.keys(newErrors);
  if (prevKeys.length !== newKeys.length) {
    changed = true;
  } else {
    for (const k of newKeys) {
      const a = prev[k] ?? [];
      const b = newErrors[k];
      if (a.length !== b.length) {
        changed = true;
        break;
      }
      for (let i = 0; i < a.length; i++) {
        if (a[i] !== b[i]) {
          changed = true;
          break;
        }
      }
      if (changed) break;
    }
  }

  if (changed) {
    prevRef.current = newErrors;
  }

  return prevRef.current;
}
