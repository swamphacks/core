import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";
import { useQuery } from "@tanstack/react-query";

import { z } from "zod";

export const ApplicationFieldsSchema = z.object({
  firstName: z.string().max(50),
  lastName: z.string().max(50),
  age: z.string(),
  phone: z.string().length(10),
  preferredEmail: z.email(),
  universityEmail: z.email(),
  country: z.string(),
  gender: z.string().optional(),
  "gender-other": z.string().optional(),
  pronouns: z.string().optional(),
  race: z.string().optional(),
  "race-other": z.string().optional(),
  orientation: z.string().optional(),
  linkedin: z.url(),
  // github: z.url(),
  ageCertification: z.boolean(),
  school: z.string(),
  level: z.string(),
  "level-other": z.string().optional(),
  year: z.string(),
  "year-other": z.string().optional(),
  graduationYear: z.string(),
  majors: z.string(),
  minors: z.string().optional(),
  experience: z.string(),
  ufHackathonExp: z.string(),
  projectExperience: z.string(),
  shirtSize: z.string(),
  diet: z.string().optional(),
  essay1: z.string(),
  essay2: z.string(),
  essay3: z.string(),
  referral: z.string(),
  pictureConsent: z.string(),
  inpersonAcknowledgement: z.string(),
  agreeToConduct: z.string(),
  infoShareAuthorization: z.string(),
  agreeToMLHEmails: z.string().optional(),
});

export type ApplicationResponse =
  operations["get-my-application"]["responses"]["200"]["content"]["application/json"];

export type ApplicationFields = z.infer<typeof ApplicationFieldsSchema>;

export type Application = Omit<ApplicationResponse, "application"> & {
  application: ApplicationFields;
};

export async function fetchApplication(userId: string): Promise<Application> {
  const result = await api.get<Application>(`application/${userId}`).json();

  const parsedApplication = ApplicationFieldsSchema.safeParse(
    JSON.parse(atob(result.application as unknown as string)),
  );
  if (!parsedApplication.success) {
    console.error("Failed to parse application data:", parsedApplication.error);
    throw new Error("Invalid application data format");
  }

  return {
    ...result,
    application: parsedApplication.data,
  };
}

export function useApplication(userId: string) {
  return useQuery({
    queryKey: ["application", userId],
    queryFn: () => fetchApplication(userId),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}
