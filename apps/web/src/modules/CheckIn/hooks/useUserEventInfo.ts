import { api } from "@/lib/ky";
import { useQuery } from "@tanstack/react-query";
import { z } from "zod";

const userSchema = z.object({
  id: z.uuid(),
  name: z.string(),
  email: z.email(),
  image: z.string().nullable().optional(),
  role: z.enum(["admin", "staff", "attendee", "applicant", "visitor"]),
  checked_in_at: z.coerce.date().nullable(),
});

async function fetchUserEventInfo(userId: string) {
  const result = await api.get(`users/userid/${userId}`).json();
  const user = userSchema.parse(result);

  return {
    user_id: user.id,
    name: user.name,
    email: user.email,
    image: user.image ?? null,
    event_role: user.role === "attendee" ? ("attendee" as const) : null,
    checked_in_at: user.checked_in_at,
  };
}

export type UserEventInfo = Awaited<ReturnType<typeof fetchUserEventInfo>>;

export function useUserEventInfo(eventId: string, userId: string | null) {
  return useQuery({
    queryKey: ["userEventInfo", eventId, userId],
    queryFn: () => {
      if (!userId) throw new Error("User ID is required");
      return fetchUserEventInfo(userId);
    },
    enabled: !!userId,
  });
}
