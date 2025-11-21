import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const ParagraphQuestion = createQuestionItem({
  type: QuestionTypes.paragraph,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.paragraph),
    showWordCount: z.boolean().default(false),
    validation: z
      .object({
        min: z.number(),
        max: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.paragraph];
    const requiredMessage = item.requiredMessage ?? error.required;

    let schema = z.string(requiredMessage);

    const countWords = (text: string) => {
      const words = text.trim().split(/\s+/);
      return words.filter((word) => word !== "").length;
    };

    if (item.isRequired) {
      schema = schema.refine((val) => countWords(val) >= 1, {
        message: requiredMessage,
      });
    }

    const { validation } = item;
    if (!validation) return schema;

    if (typeof validation.max === "number" && validation.max !== undefined) {
      schema = schema.refine((val) => countWords(val) <= validation.max!, {
        message: error.tooLong,
      });
    }

    if (typeof validation.min === "number" && validation.min !== undefined) {
      schema = schema.refine((val) => countWords(val) >= validation.min!, {
        message: error.tooShort,
      });
    }

    return schema;
  },
});
