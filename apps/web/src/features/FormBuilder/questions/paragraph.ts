import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const ParagraphQuestion = createQuestionItem({
  type: QuestionTypes.paragraph,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.paragraph),
    validation: z
      .object({
        minLength: z.number(),
        maxLength: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.paragraph];
    const requiredMessage = item.requiredMessage ?? error.required;

    let schema = z.string(requiredMessage);

    if (item.isRequired) {
      schema = schema.min(1, requiredMessage);
    }

    const { validation } = item;
    if (!validation) return schema;

    if (typeof validation.maxLength === "number") {
      schema = schema.max(validation.maxLength, error.tooLong);
    }

    if (typeof validation.minLength === "number") {
      schema = schema.min(validation.minLength, error.tooShort);
    }

    return schema;
  },
});
