import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const DateQuestion = createQuestionItem({
  type: QuestionTypes.date,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.date),
    validation: z
      .object({
        min: z.date(),
        max: z.date(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.date];
    const requiredMessage = item.requiredMessage ?? error.required;

    const schema = z.date(requiredMessage);

    const { validation } = item;
    if (!validation) return schema;

    if (validation.min) {
      schema.min(validation.min);
    }

    if (validation.max) {
      schema.max(validation.max);
    }

    return schema;
  },
});
