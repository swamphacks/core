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
    const { validation } = item;

    const schema = z.string().check((ctx) => {
      const urlCheck = z.url(error.invalidURL).safeParse(ctx.value);
      if (!urlCheck.success) {
        ctx.issues.push({
          code: "custom",
          message: error.invalidURL,
          input: ctx.value,
        });
        return;
      }

      if (item.regex) {
        const re = new RegExp(item.regex);
        if (!re.test(ctx.value)) {
          ctx.issues.push({
            code: "custom",
            message: "Invalid value",
            input: ctx.value,
          });
          return;
        }
      }

      if (validation && typeof validation.maxLength === "number") {
        if (ctx.value.length > validation.maxLength) {
          ctx.issues.push({
            code: "custom",
            message: error.tooLong,
            input: ctx.value,
          });
          return;
        }
      }
    });

    return schema;
  },
});
