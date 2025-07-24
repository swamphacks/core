import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const DateQuestion = createQuestionItem({
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
    const schema = z.date();

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
