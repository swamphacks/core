import z from "zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import { queryKey as authQueryKey } from "@/lib/auth/hooks/useUser";

export const settingsFieldsSchema = z.object({
  name: z.string().min(1, "Name is required"),
  preferredEmail: z.preprocess(
    (val: string) => (val === "" ? undefined : val),
    z.email("Email address is invalid").optional(),
  ),
});

export function useSettingsActions() {
  const queryClient = useQueryClient();

  const updateAccountInfo = useMutation({
    mutationFn: (data: z.infer<typeof settingsFieldsSchema>) => {
      return api
        .patch("users/me", {
          json: {
            name: data.name,
            preferred_email: data.preferredEmail,
          },
        })
        .json();
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: authQueryKey });
    },
  });

  return { updateAccountInfo };
}
