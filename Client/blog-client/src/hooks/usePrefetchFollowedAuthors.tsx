import { getFollowedAuthorsRequest } from "@/api/authorApi";
import { useQueryClient } from "@tanstack/react-query";
import { useCallback } from "react";

export function usePrefetchFollowedAuthor() {
  const queryClient = useQueryClient();

  return useCallback(() => {
    return queryClient.prefetchQuery({
      queryKey: ["followed_authors"],
      queryFn: () => getFollowedAuthorsRequest(),
      staleTime: 1000 * 60 * 5,
    });
  }, [queryClient]);
}
