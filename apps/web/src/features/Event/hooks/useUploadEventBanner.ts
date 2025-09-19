import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Event } from "../schemas/event";
import { api } from "@/lib/ky";
import z from "zod";

//TODO: change this into a openapi spec
const BannerUploadResponseSchema = z.object({
  banner_url: z.url(),
});

export const useUploadEventBanner = (eventId: string) => {
  const queryClient = useQueryClient();

  const eventQueryKey = ["event", eventId] as const;

  const uploadBanner = async (bannerFile: File) => {
    const formData = new FormData();
    formData.append("image", bannerFile);

    const result = await api
      .post(`events/${eventId}/banner`, {
        body: formData,
      })
      .json();

    return BannerUploadResponseSchema.parse(result);
  };

  return useMutation({
    mutationFn: uploadBanner,
    onSuccess: (data) => {
      queryClient.setQueryData<Event>(eventQueryKey, (old) => {
        if (!old) return;

        return {
          ...old,
          banner: data.banner_url,
        };
      });

      // Also return the new banner URL
      return data.banner_url;
    },
  });
};
