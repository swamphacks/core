import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const CheckboxQuestion = createQuestionItem({
  type: QuestionTypes.checkbox,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.checkbox),
    options: z.array(
      z.object({
        label: z.string(),
        value: z
          .string()
          .or(z.number())
          .transform((val) => (typeof val === "number" ? val.toString() : val)), // convert value to string
      }),
    ),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.checkbox];
    const requiredMessage = item.requiredMessage ?? error.required;

    let schema = z.array(z.string(), requiredMessage); // array of string value of selections

    if (item.isRequired) {
      schema = schema.min(1, requiredMessage);
    }

    return schema;
  },
});
