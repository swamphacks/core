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
    if (!validation) return z.array(schema, "File(s) must be uploaded.");

    if (validation.validMimeTypes) {
      schema = z.file().mime(validation.validMimeTypes, "Invalid file type");
    }

    if (validation.invalidMimeTypes) {
      schema = schema.refine(
        (file) => !validation.invalidMimeTypes?.includes(file.type),
        "Invalid file type.",
      );
    }

    // default 5MB
    const { maxSize = 5 * 1024 * 1024 } = validation;

    schema = schema.refine((file) => {
      if (validation.minSize) {
        return validation.minSize <= file.size && file.size <= maxSize;
      }

      return file.size <= maxSize;
    }, "File size is not within range.");

    const flattenedArraySchema = z.any().check((ctx) => {
      if (!Array.isArray(ctx.value)) {
        ctx.issues.push({
          code: "custom",
          message: "File(s) must be uploaded.",
          input: ctx.value,
        });
        return;
      }

      for (const file of ctx.value) {
        const result = schema.safeParse(file);
        if (!result.success) {
          for (const issue of result.error.issues) {
            ctx.issues.push({
              code: "custom",
              message: issue.message,
              input: ctx.value,
            });
          }
        }
      }
    });

    return flattenedArraySchema;
  },
});
