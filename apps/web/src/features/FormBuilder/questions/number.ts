import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const NumberQuestion = createQuestionItem({
  type: QuestionTypes.number,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.number),
    validation: z
      .object({
        min: z.number(),
        max: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.number];
    const requiredMessage = item.requiredMessage ?? error.required;

    let schema = z.string(requiredMessage);

    if (item.isRequired) {
      schema = schema.min(1, requiredMessage);
    }

    const newSchema = schema.transform((val) => parseInt(val));

    const { validation } = item;
    if (!validation) return newSchema;

    let numberSchema = z.number();

    if (typeof validation.min === "number") {
      numberSchema = numberSchema.min(validation.min, error.tooLow);
    }

    if (typeof validation.max === "number") {
      numberSchema = numberSchema.max(validation.max, error.tooHigh);
    }

    return newSchema.pipe(numberSchema);
  },
});
