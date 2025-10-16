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
    regex: z.string().optional(),
    validation: z
      .object({
        maxLength: z.number(),
        email: z.boolean(),
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

    if (item.regex) {
      schema = schema.regex(
        new RegExp(item.regex as unknown as RegExp),
        "Invalid value",
      );
    }

    const { validation } = item;
    if (!validation) return schema;

    if (validation.email) {
      let schema = z.email("Invalid email");

      if (item.regex) {
        schema = schema.regex(
          new RegExp(item.regex as unknown as RegExp),
          "Invalid value",
        );
      }

      if (typeof validation.maxLength === "number") {
        schema = schema.max(validation.maxLength, error.tooLong);
      }

      return schema;
    }

    if (typeof validation.maxLength === "number") {
      schema = schema.max(validation.maxLength, error.tooLong);
    }

    return schema;
  },
});
