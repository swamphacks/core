import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/ky";
import type { Dashboard } from "../types";

interface DashboardFilters {
  hackathonId?: string;
  search?: string;
  status?: string;
  redeemableId?: string;
}

export function useCheckInDashboard(
  {
    hackathonId = "",
    search = "",
    status = "",
    redeemableId = "",
  }: DashboardFilters = {},
  options: { enabled?: boolean } = {},
) {
  return useQuery({
    queryKey: ["check-in-dashboard", hackathonId, search, status, redeemableId],
    queryFn: ({ signal }) =>
      api
        .get("application/check-in-dashboard", {
          searchParams: { hackathonId, search, status, redeemableId },
          signal,
        })
        .json<Dashboard>(),
    refetchInterval: 15000,
    enabled: options.enabled,
  });
}
