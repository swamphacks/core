import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";

import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const MultiSelectQuestion = createQuestionItem({
  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.multiselect),
    options: z
      .array(
        z.object({
          label: z.string(),
          value: z
            .string()
            .or(z.number())
            .transform((val) =>
              typeof val === "number" ? val.toString() : val,
            ), // convert value to string
        }),
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

  extractValidationSchemaFromItem: (item) => {
    let schema = z.array(z.string(), "Pick an item or more from the list."); // array of string value of selections

    if (item.required) {
      schema = schema.min(1, "pcik an item");
    }

    return schema;
  },
});
