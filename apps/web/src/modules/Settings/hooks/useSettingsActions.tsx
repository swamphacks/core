import z from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import { useUserQueryKey } from "@/lib/auth/hooks/useUser";

export const settingsFieldsSchema = z.object({
  name: z.string().min(1, "Name is required"),
  preferredEmail: z.preprocess((val: string) => {
    if (typeof val !== "string") return undefined;
    if (val.trim() === "") return undefined;
    return val;
  }, z.email("Email address is invalid").optional()),
});

export function useSettingsActions() {
  const queryClient = useQueryClient();

  const updateAccountInfo = useMutation({
    mutationFn: (data: z.infer<typeof settingsFieldsSchema>) => {
      return api
        .patch("users/me", {
          json: {
            name: data.name,
            preferredEmail: data.preferredEmail,
          },
        })
        .json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: useUserQueryKey });
    },
  });

  const updateEmailConsent = useMutation({
    mutationFn: (selected: boolean) => {
      return api
        .patch("users/me/email-consent", {
          json: { emailConsent: selected },
        })
        .json();
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: useUserQueryKey });
    },
  });

  return { updateAccountInfo, updateEmailConsent };
}
