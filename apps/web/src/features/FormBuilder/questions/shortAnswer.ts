import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const ShortAnswerQuestion = createQuestionItem({
  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.shortAnswer),
    validation: z
      .object({
        maxLength: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    let schema = z.string();

    const { validation } = item;
    if (!validation) return schema;

    if (typeof validation.maxLength === "number") {
      schema = schema.max(
        validation.maxLength,
        "Value exceeded character limit.",
      );
    }

    return schema;
  },
});
