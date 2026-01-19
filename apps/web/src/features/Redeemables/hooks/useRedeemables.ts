import { api } from "@/lib/ky";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

export interface Redeemable {
  id: string;
  event_id: string;
  name: string;
  total_stock: number;
  max_user_amount: number;
  created_at: string;
  updated_at: string;
  total_redeemed: number | string;
}

export async function fetchRedeemables(eventId: string): Promise<Redeemable[]> {
  const result = await api
    .get<Redeemable[]>(`events/${eventId}/redeemables`)
    .json();
  return result;
}

export function useRedeemables(eventId: string) {
  return useQuery({
    queryKey: ["redeemables", eventId],
    queryFn: () => fetchRedeemables(eventId),
  });
}

interface CreateRedeemableRequest {
  name: string;
  amount: number;
  max_user_amount: number;
}

export function useCreateRedeemable(eventId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: CreateRedeemableRequest) => {
      // The backend returns a Redeemable struct (which has 'amount' field)
      // but we invalidate queries so the GET response (with 'total_stock') will be fetched
      await api.post(`events/${eventId}/redeemables`, {
        json: data,
      });
    },
    onSuccess: () => {
      // Invalidate to refetch with the correct format from GET endpoint
      queryClient.invalidateQueries({ queryKey: ["redeemables", eventId] });
    },
  });
}

interface GetUserByRFIDResponse {
  user_id: string;
}

// Convert RFID to UserID
export async function getUserByRFID(
  eventId: string,
  rfid: string,
): Promise<GetUserByRFIDResponse> {
  const result = await api
    .get<GetUserByRFIDResponse>(`events/${eventId}/users/by-rfid/${rfid}`)
    .json();
  return result;
}

export function useRedeemRedeemable(eventId: string, redeemableId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (userId: string) => {
      await api.post(
        `events/${eventId}/redeemables/${redeemableId}/users/${userId}`,
      );
    },
    onSuccess: () => {
      // Invalidate to refetch updated redeemable data
      queryClient.invalidateQueries({ queryKey: ["redeemables", eventId] });
    },
  });
}
