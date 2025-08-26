import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";

type NullableEventRoleRow =
  operations["get-user-event-role"]["responses"]["200"]["content"]["application/json"];

export async function getUserEventRole(
  eventId: string,
): Promise<NullableEventRoleRow> {
  const result = await api
    .get<NullableEventRoleRow>(`events/${eventId}/role`)
    .json();

  return result;
}
