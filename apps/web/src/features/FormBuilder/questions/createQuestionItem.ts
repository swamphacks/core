import type { BaseFormQuestionItemSchema } from "../formSchema";
import type { QuestionTypes } from "../types";
import { BaseQuestion } from "./baseQuestion";
import z from "zod";

type AnyQuestionZodSchema =
  (typeof BaseFormQuestionItemSchema.def.options)[number];

export const questionTypeItemMap = {} as Record<
  keyof typeof QuestionTypes,
  QuestionItem<AnyQuestionZodSchema>
>;

export interface QuestionItem<
  T extends typeof BaseQuestion = typeof BaseQuestion,
  TReturn = z.ZodAny,
> {
  type: keyof typeof QuestionTypes;
  schema: T;
  extractValidationSchemaFromItem: (
    item: Omit<z.infer<T>, "required"> & { isRequired: boolean },
  ) => TReturn;
}

// This function allows us to create stricter typing for question items without having to manually pass in generic types to the QuestionItem interface.
export function createQuestionItem<
  TSchema extends QuestionItem["schema"],
  TReturn,
>(item: QuestionItem<TSchema, TReturn>): QuestionItem<TSchema, TReturn> {
  // technically, this function is not supposed to know the type of `item` here because its job is to infer that type
  // but we still want to assert that we know its type already so we can store it in the map
  questionTypeItemMap[item.type] =
    item as unknown as QuestionItem<AnyQuestionZodSchema>;
  return item;
}
