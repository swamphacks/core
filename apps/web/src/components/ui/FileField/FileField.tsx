// https://github.com/adobe/react-spectrum/discussions/5486#discussioncomment-10264448 (with some modifications)
// The FileTrigger component doesn't support form submission yet, hence we use a custom file field component here.

import { useFormValidation } from "@react-aria/form";
import {
  FormValidationContext,
  useFormValidationState,
} from "@react-stately/form";
import {
  createContext,
  type JSX,
  type PropsWithChildren,
  useContext,
  useId,
  useMemo,
  useRef,
} from "react";
import { FieldErrorContext, InputContext } from "react-aria-components";
import { Description, FieldError, Label } from "@/components/ui/Field";
import { FileInput } from "./FileInput";
import type { ValidationError } from "@tanstack/react-form";

export interface FileFieldProps {
  name: string;
  isRequired?: boolean | undefined;
  accept?: string | undefined;
  multiple?: boolean | undefined;
  description?: string;
  label?: string;
  onChange?: (files: File[]) => void;
  maxSize?: number;
  errorMessage?: string;
  validationBehavior?: "aria" | "native";
}

// This context allows FileField to accept arbitrary props so that the FileInput component can use without affecting
// the props of the underlying input component
export const CustomPropsContext = createContext<{
  onChange?: (args: any) => void;
  maxSize?: number;
  error?: ValidationError;
  resetValidation: () => void;
}>(undefined!);

export const FileField = (props: FileFieldProps) => {
  return (
    <FileFieldWrapper {...props}>
      <FileInput />
      <div>
        <Description>{props.description}</Description>
      </div>
      <FieldError />
    </FileFieldWrapper>
  );
};

export const FileFieldWrapper = ({
  children,
  name,
  isRequired,
  accept,
  multiple,
  label,
  onChange,
  maxSize,
  errorMessage,
  validationBehavior = "aria",
}: PropsWithChildren<FileFieldProps>): JSX.Element => {
  const inputRef = useRef<HTMLInputElement>(null);
  const id = useId();

  const formValidationState = useFormValidationState({
    value: undefined,
    validationBehavior,
  });

  useFormValidation({ validationBehavior }, formValidationState, inputRef);

  const errors = useContext(FormValidationContext);
  const formError = errors?.[name] ?? [];

  // Prioritize error passed in by errorMessage prop and from Form component using validationErrors prop over native errors
  const isInvalid =
    !!errorMessage ||
    formError.length > 0 ||
    formValidationState.displayValidation.isInvalid;

  const getError = () => {
    if (errorMessage) {
      return {
        ...formValidationState.displayValidation,
        isInvalid,
        validationErrors: [errorMessage],
      };
    }

    if (formError.length > 0) {
      return {
        ...formValidationState.displayValidation,
        isInvalid,
        validationErrors: [...formError],
      };
    }

    return formValidationState.displayValidation;
  };

  const inputProps = useMemo(
    () => ({
      id,
      name,
      // If validationBehavior is "aria", we expect `required` to be validated using a zod schema instead of the native behavior
      required: validationBehavior === "aria" ? false : isRequired,
      accept,
      multiple,
      ref: inputRef,
      "aria-invalid": isInvalid,
    }),
    [name, isRequired, accept, multiple, formValidationState],
  );

  return (
    <div
      data-required={!!isRequired || undefined}
      data-invalid={isInvalid || undefined}
      className="flex flex-col gap-1"
    >
      <Label htmlFor={id} isRequired={isRequired}>
        {label}
      </Label>

      <CustomPropsContext.Provider
        value={{
          onChange,
          maxSize,
          resetValidation: formValidationState.resetValidation,
        }}
      >
        <InputContext.Provider value={inputProps}>
          <FieldErrorContext.Provider value={getError()}>
            {children}
          </FieldErrorContext.Provider>
        </InputContext.Provider>
      </CustomPropsContext.Provider>
    </div>
  );
};
