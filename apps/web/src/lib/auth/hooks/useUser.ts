import { useQuery } from "@tanstack/react-query";
import { _getUser } from "../services/user";

export const queryKey = ["auth", "me"] as const;

export function _useUser() {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  return useQuery({
    queryKey,
    queryFn: async () => await _getUser(),
    refetchOnWindowFocus: false,
    refetchOnMount: "always",
    staleTime: 1000 * 60 * 10, // 10 minutes
    retry: false,
  });
}
