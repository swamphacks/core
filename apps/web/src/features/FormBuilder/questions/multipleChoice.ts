import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const MultipleChoiceQuestion = createQuestionItem({
  type: QuestionTypes.multipleChoice,

  // Does multiple choice question have any validation besides being required?
  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.multipleChoice),
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

  extractValidationSchemaFromItem: (item) =>
    z.string(
      item.requiredMessage ??
        errorMessage[QuestionTypes.multipleChoice].required,
    ),
});
