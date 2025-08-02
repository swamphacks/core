import { FormItemTypes } from "@/features/FormBuilder/types";
import z from "zod";

export const BaseQuestion = z.object({
  type: z.literal(FormItemTypes.question),
  id: z.string().min(1, "Question must have a unique string id."),
  name: z
    .string()
    .min(1, "Question must have a unique name for form submissions."),
  label: z.string().min(1, "Question must have a label."),
  required: z.boolean().optional().default(false),
  placeholder: z.string().optional(),
  description: z.string().optional(),
});
