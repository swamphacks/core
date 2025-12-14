import { api } from "@/lib/ky";
import type { paths } from "@/lib/openapi/schema";

type NullableEventRoleRow =
  | paths["/events/{eventId}/role"]["get"]["responses"]["200"]["content"]["application/json"]
  | undefined
  | null;

export async function getUserEventRole(
  eventId: string,
): Promise<NullableEventRoleRow> {
  const result = await api
    .get<NullableEventRoleRow>(`events/${eventId}/role`)
    .json();

  return result;
}
