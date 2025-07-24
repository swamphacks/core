import { createQuestionItem } from "@/features/FormBuilder/questions/createQuestionItem";
import { QuestionTypes } from "@/features/FormBuilder/types";
import z from "zod";
import { BaseQuestion } from "./baseQuestion";

export const UploadQuestion = createQuestionItem({
  schema: BaseQuestion.extend({
    questionType: z.literal(QuestionTypes.upload),
    validation: z
      .object({
        // Size are in megabytes
        minSize: z.number(),
        maxSize: z.number(),
        validMimeTypes: z.array(z.string()),
        invalidMimeTypes: z.array(z.string()),
      })
      .partial()
      .refine(
        (data) => !(data.validMimeTypes && data.invalidMimeTypes), // both cannot be defined
        {
          message:
            "You can define either validMimeTypes or invalidMimeTypes, not both.",
        },
      )
      .optional(),
  }),

  extractValidationSchemaFromItem: (item) => {
    let schema = z.file();

    const { validation } = item;
    if (!validation) return z.array(schema);

    const { maxSize = 10 } = validation;

    if (Array.isArray(validation.validMimeTypes)) {
      schema = schema.mime(validation.validMimeTypes);
    }

    if (Array.isArray(validation.invalidMimeTypes)) {
      schema = schema.refine(
        (file) => !validation.invalidMimeTypes?.includes(file.type),
        "Invalid file type.",
      );
    }

    schema = schema.refine((file) => {
      if (validation.minSize) {
        return validation.minSize <= file.size && file.size <= maxSize;
      }

      return file.size <= maxSize;
    }, "File size is not within range.");

    return z.array(schema);
  },
});
