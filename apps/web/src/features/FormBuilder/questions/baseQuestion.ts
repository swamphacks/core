import { FormItemTypes } from "@/features/FormBuilder/types";
import z from "zod";

export const BaseQuestion = z.object({
  type: z
    .literal(FormItemTypes.question)
    .optional()
    .default(FormItemTypes.question),
  // id: z.string().min(1, "Question must have a unique string id."),
  name: z
    .string()
    .min(1, "Question must have a unique name for form submissions."),
  label: z.string().optional(),
  renderLabelAsHTML: z.boolean().default(false),
  required: z.boolean().optional().default(false), // this prop will be transformed to `isRequired` after form parsing
  requiredMessage: z.string().optional(),
  placeholder: z.string().optional(),
  description: z.string().optional(),
});
