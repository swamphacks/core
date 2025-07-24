import { BaseQuestion } from "./baseQuestion";
import z from "zod";

export interface QuestionItem<
  T extends typeof BaseQuestion = typeof BaseQuestion,
  TReturn = z.ZodAny,
> {
  schema: T;
  extractValidationSchemaFromItem: (item: z.infer<T>) => TReturn;
}

// This function allows us to create stricter typing for question items without having to manually pass in generic types to the QuestionItem interface.
export function createQuestionItem<
  TSchema extends QuestionItem["schema"],
  TReturn,
>(item: QuestionItem<TSchema, TReturn>): QuestionItem<TSchema, TReturn> {
  return item;
}
