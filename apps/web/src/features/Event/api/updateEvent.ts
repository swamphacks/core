import { api } from "@/lib/ky";
import type { Event } from "../schemas/event";

// TODO: Change this from schemas event to the openapi event

export async function updateEventById(
  eventId: string,
  data: Partial<Event>,
): Promise<Event> {
  const result = await api
    .patch<Event>(`events/${eventId}`, { json: data })
    .json();

  return result;
}
