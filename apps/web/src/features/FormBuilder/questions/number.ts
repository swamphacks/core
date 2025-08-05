import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const NumberQuestion = createQuestionItem({
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
    let schema = z.string("Fill out this field.");

    if (item.required) {
      schema = schema.min(1, "Fill out this field");
    }

    const newSchema = schema.transform((val) => parseInt(val));

    const { validation } = item;
    if (!validation) return newSchema;

    let numberSchema = z.number();

    if (typeof validation.min === "number") {
      numberSchema = numberSchema.min(validation.min, "Number is too small.");
    }

    if (typeof validation.max === "number") {
      numberSchema = numberSchema.max(validation.max, "Number is too big.");
    }

    return newSchema.pipe(numberSchema);
  },
});
