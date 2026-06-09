import { api } from "@/lib/ky";
import type { AssignedApplications } from "@/lib/openapi/types";
import { useQuery } from "@tanstack/react-query";

export const assignedApplicationsQueryKey = ["assignedApplications"];

export const fetchAssignedApplications = async () => {
  const data = await api
    .get<AssignedApplications>(`application/assigned`)
    .json();

  return data;
};

export function useAssignedApplications() {
  return useQuery({
    queryKey: assignedApplicationsQueryKey,
    queryFn: () => fetchAssignedApplications(),
    staleTime: 1000 * 60 * 15, // 15 minutes,
  });
}
