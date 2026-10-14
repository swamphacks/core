import { api } from "@/lib/ky";
import type { components } from "@/lib/openapi/schema";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { HTTPError } from "ky";

export type Redeemable = components["schemas"]["GetRedeemablesRow"];
export type CreateRedeemableRequest =
  components["schemas"]["CreateRedeemableRequest"];
export type UpdateRedeemableRequest =
  components["schemas"]["UpdateRedeemableRequest"];
export type RedeemableType = CreateRedeemableRequest["type"];

// The NFC app lists items by type, one tab per type. Swag and drink have no
// tab yet, so items of those types can't be scanned until the app adds them.
export const REDEEMABLE_TYPES: {
  id: RedeemableType;
  name: string;
  appTab: string | null;
}[] = [
  { id: "meal", name: "Meal", appTab: "Meals" },
  { id: "tshirt", name: "T-shirt", appTab: "T-shirts" },
  { id: "swag", name: "Swag", appTab: null },
  { id: "drink", name: "Drink", appTab: null },
];

export function redeemableTypeInfo(type: string) {
  return REDEEMABLE_TYPES.find((t) => t.id === type);
}

export const redeemablesQueryKey = ["redeemables"];

export function useRedeemables() {
  return useQuery({
    queryKey: redeemablesQueryKey,
    queryFn: async () =>
      (await api.get("redeemables").json<Redeemable[] | null>()) ?? [],
    // Keeps the redeemed counts current while scanning is happening.
    refetchInterval: 15000,
  });
}

export function useCreateRedeemable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: CreateRedeemableRequest) => {
      await api.post("redeemables", { json: body });
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: redeemablesQueryKey }),
  });
}

export function useUpdateRedeemable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({
      id,
      body,
    }: {
      id: string;
      body: UpdateRedeemableRequest;
    }) => {
      await api.patch(`redeemables/${id}`, { json: body });
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: redeemablesQueryKey }),
  });
}

export function useDeleteRedeemable() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await api.delete(`redeemables/${id}`);
    },
    onSuccess: () =>
      queryClient.invalidateQueries({ queryKey: redeemablesQueryKey }),
  });
}

// The API returns huma error bodies: { detail, errors?: [{ message, location }] }.
export async function apiErrorMessage(error: unknown, fallback: string) {
  if (!(error instanceof HTTPError)) return fallback;

  try {
    const body = (await error.response.json()) as {
      detail?: string;
      errors?: { message?: string; location?: string }[];
    };
    const fieldErrors = body.errors
      ?.map((e) => e.message)
      .filter(Boolean)
      .join("; ");
    return fieldErrors || body.detail || fallback;
  } catch {
    return fallback;
  }
}
