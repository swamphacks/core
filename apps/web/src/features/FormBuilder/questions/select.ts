import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";
import { errorMessage } from "../errorMessage";

export const SelectQuestion = createQuestionItem({
  type: QuestionTypes.select,

  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.select),
    searchable: z.boolean().default(false),
    hasOther: z.boolean().default(false),
    options: z
      .array(
        z
          .object({
            label: z.string(),
            value: z
              .string()
              .or(z.number())
              .transform((val) =>
                typeof val === "number" ? val.toString() : val,
              ), // convert value to string
          })
          .transform((data) => ({
            id: data.value,
            name: data.label,
          })),
      )
      .or(
        z
          .object({
            data: z.enum(
              ["schools", "majors", "minors", "year", "countries"],
              "Invalid data option for select question.",
            ),
            min: z.number().optional(),
            max: z.number().optional(),
          })
          .refine((val) => {
            if (val.data === "year") {
              return val.min && val.max;
            }

            return true;
          }),
      ),
  }),

  extractValidationSchemaFromItem: (item) =>
    z.string(
      item.requiredMessage ?? errorMessage[QuestionTypes.select].required,
    ), // string value of selection
});
