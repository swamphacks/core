import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";
import { textFieldIcons } from "@/features/FormBuilder/icons";

export const URLQuestion = createQuestionItem({
  type: QuestionTypes.url,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.url),
    iconName: z
      .enum(Object.keys(textFieldIcons) as Array<keyof typeof textFieldIcons>)
      .optional(),
    regex: z.string().optional(),
    validation: z
      .object({
        maxLength: z.number(),
      })
      .partial()
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    const error = errorMessage[QuestionTypes.url];
    // const requiredMessage = item.requiredMessage ?? error.required;

    let schema = z.url(error.invalidURL);

    // if (item.isRequired) {
    //   schema = schema.min(1, requiredMessage);
    // }

    if (item.regex) {
      schema = schema.regex(
        new RegExp(item.regex as unknown as RegExp),
        "Invalid value",
      );
    }

    const { validation } = item;
    if (!validation) return schema;

    if (typeof validation.maxLength === "number") {
      schema = schema.max(validation.maxLength, error.tooLong);
    }

    return schema;
  },
});
