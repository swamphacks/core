// https://github.com/adobe/react-spectrum/discussions/5486#discussioncomment-10264448 (with some modifications)
// The FileTrigger component doesn't support form submission yet, hence we use a custom file field component here.

import { useFormValidation } from "@react-aria/form";
import { useFormValidationState } from "@react-stately/form";
import {
  createContext,
  type JSX,
  type PropsWithChildren,
  useId,
  useMemo,
  useRef,
} from "react";
import { FieldErrorContext, InputContext } from "react-aria-components";
import { Description, FieldError, Label } from "@/components/ui/Field";
import { FileInput } from "./FileInput";

export interface FileFieldProps {
  name?: string | undefined;
  isRequired?: boolean | undefined;
  accept?: string | undefined;
  multiple?: boolean | undefined;
  description?: string;
  label?: string;
  onChange?: (files: File[]) => void;
}

// This context allows FileField to accept arbitrary props so that the FileInput component can use without affecting
// the props of the underlying input component
export const CustomPropsContext = createContext<
  | {
      onChange?: (args: any) => void;
    }
  | undefined
>(undefined);

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
}: PropsWithChildren<FileFieldProps>): JSX.Element => {
  const inputRef = useRef<HTMLInputElement>(null);
  const id = useId();

  const formValidationState = useFormValidationState({
    value: undefined,
    validationBehavior: "native",
  });

  useFormValidation(
    { validationBehavior: "native" },
    formValidationState,
    inputRef,
  );

  const inputProps = useMemo(
    () => ({
      id,
      name,
      required: isRequired,
      accept,
      multiple,
      ref: inputRef,
      "aria-invalid": formValidationState.displayValidation.isInvalid,
    }),
    [name, isRequired, accept, multiple, formValidationState],
  );

  return (
    <div
      data-required={!!isRequired || undefined}
      data-invalid={
        formValidationState.displayValidation.isInvalid || undefined
      }
      className="flex flex-col gap-1"
    >
      <Label htmlFor={id} isRequired={isRequired}>
        {label}
      </Label>

      <CustomPropsContext.Provider value={{ onChange }}>
        <InputContext.Provider value={inputProps}>
          <FieldErrorContext.Provider
            value={formValidationState.displayValidation}
          >
            {children}
          </FieldErrorContext.Provider>
        </InputContext.Provider>
      </CustomPropsContext.Provider>
    </div>
  );
};
