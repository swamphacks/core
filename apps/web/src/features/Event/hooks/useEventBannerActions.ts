import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Event } from "../schemas/event";
import { api } from "@/lib/ky";
import z from "zod";

// Response schema for banner upload
const BannerUploadResponseSchema = z.object({
  banner_url: z.url(),
});

export const useEventBannerActions = (eventId: string) => {
  const queryClient = useQueryClient();
  const eventQueryKey = ["event", eventId] as const;

  // Upload banner
  const upload = useMutation({
    mutationFn: async (bannerFile: File) => {
      const formData = new FormData();
      formData.append("image", bannerFile);

      const result = await api
        .post(`events/${eventId}/banner`, { body: formData })
        .json();
      return BannerUploadResponseSchema.parse(result);
    },
    onSuccess: (data) => {
      queryClient.setQueryData<Event>(eventQueryKey, (old) => {
        if (!old) return old;

        return {
          ...old,
          banner: data.banner_url,
        };
      });
    },
  });

  // Delete banner
  const remove = useMutation({
    mutationFn: async () => {
      await api.delete(`events/${eventId}/banner`).json();

      return null;
    },
    onSuccess: () => {
      queryClient.setQueryData<Event>(eventQueryKey, (old) => {
        if (!old) return old;

        return {
          ...old,
          banner: null,
        };
      });
    },
  });

  return { upload, remove };
};
