import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const SelectQuestion = createQuestionItem({
  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.select),
    searchable: z.boolean().default(false),
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
        z.object({
          data: z.enum(
            ["schools", "majors", "minors"],
            "Invalid data option for select question.",
          ),
        }),
      ),
  }),

  extractValidationSchemaFromItem: () => z.string("Pick an item."), // string value of selection
});
