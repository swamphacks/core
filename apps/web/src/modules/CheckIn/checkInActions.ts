import { HTTPError } from "ky";
import { api } from "@/lib/ky";

// Both return null on success, or a message to show staff.
async function post(path: string, userId: string, fallback: string) {
  try {
    await api.post(path, { json: { userId }, retry: 0 });
    return null;
  } catch (error) {
    if (error instanceof HTTPError) {
      const response = await error.response
        .json<{ detail?: string }>()
        .catch(() => null);
      return response?.detail ?? fallback;
    }
    return fallback;
  }
}

export function checkInWithoutNfc(userId: string) {
  return post("application/check-in", userId, "Unable to check in hacker.");
}

export function undoCheckIn(userId: string) {
  return post("application/check-in/undo", userId, "Unable to undo check-in.");
}
