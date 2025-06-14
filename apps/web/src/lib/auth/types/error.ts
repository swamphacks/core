import { z } from "zod";

// Error handling types for API responses
export const errorSchema = z.object({
  error: z.string(),
  message: z.string(),
});

export type ErrorResponse = z.infer<typeof errorSchema>;
