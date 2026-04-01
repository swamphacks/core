import { api } from "@/lib/ky";
import type { operations } from "@/lib/openapi/schema";

type Event =
  operations["get-single-event"]["responses"]["201"]["content"]["application/json"];

export async function getEventById(eventId: string): Promise<Event> {
  const result = await api.get<Event>(`events/${eventId}`).json();

  return result;
}
