import {
  CheckboxQuestion,
  DateQuestion,
  MultipleChoiceQuestion,
  MultiSelectQuestion,
  NumberQuestion,
  ParagraphQuestion,
  SelectQuestion,
  ShortAnswerQuestion,
  UploadQuestion,
} from "@/features/FormBuilder/questions";
import { FormItemTypes } from "@/features/FormBuilder/types";
import { z } from "zod";

export const FormQuestionItemSchema = z.discriminatedUnion("questionType", [
  ShortAnswerQuestion.schema,
  ParagraphQuestion.schema,
  NumberQuestion.schema,
  MultipleChoiceQuestion.schema,
  CheckboxQuestion.schema,
  SelectQuestion.schema,
  MultiSelectQuestion.schema,
  UploadQuestion.schema,
  DateQuestion.schema,
]);

export const FormLayoutItemSchema = z.object({
  type: z.literal(FormItemTypes.layout),
  label: z.string().optional(),
  content: z
    .array(FormQuestionItemSchema, "A layout item can only contain questions.")
    .max(2, "A layout form item cannot have more than 2 questions."),
});

export const FormSectionItemSchema = z.object({
  type: z.literal(FormItemTypes.section),
  label: z.string().optional(),
  content: z.array(
    z.union([FormLayoutItemSchema, FormQuestionItemSchema]),
    "A section cannot contain nested sections.",
  ),
});

export const FormItemSchema = z.discriminatedUnion("type", [
  FormSectionItemSchema,
  FormLayoutItemSchema,
  FormQuestionItemSchema,
]);

export const FormMetadataSchema = z.object({
  title: z.string().min(1, "Form must have a title."),
  description: z.string().optional(),
});

export const FormSchema = z.object({
  metadata: FormMetadataSchema,
  content: z.array(FormItemSchema),
});

/* eslint-disable-next-line @typescript-eslint/no-explicit-any */
export type FormObject = Record<string, any>;
export type FormItemSchemaType = z.infer<typeof FormItemSchema>;
export type FormQuestionItemSchemaType = z.infer<typeof FormQuestionItemSchema>;
