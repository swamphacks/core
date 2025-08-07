import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { textFieldIcons } from "@/features/FormBuilder/icons";
import { errorMessage } from "../errorMessage";

export const ShortAnswerQuestion = createQuestionItem({
  type: QuestionTypes.shortAnswer,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.shortAnswer),
    iconName: z
      .enum(Object.keys(textFieldIcons) as Array<keyof typeof textFieldIcons>)
      .optional(),
    validation: z
      .object({
        maxLength: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.shortAnswer];
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

    return schema;
  },
});
