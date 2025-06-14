import { useQuery } from "@tanstack/react-query";
import { _getUser } from "../services/session";

const queryKey = ["auth", "me"] as const;

// Can't have an underscore to mark as internal, maybe fix?
export function useUser() {
  const getUser = async () => {
    const response = await _getUser();

    //TODO: Figure out the best way to handle error here? Not sure if we want to throw it.
    return response;
  };

  return useQuery({
    queryKey,
    queryFn: getUser,
    refetchOnWindowFocus: false,
    refetchOnMount: "always",
    staleTime: 1000 * 60 * 10, // 10 minutes
    retry: false, // Disable retries for this query
  });
}
