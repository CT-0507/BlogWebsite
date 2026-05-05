import { getReplies } from "@/api/blogApi";
import { useQuery } from "@tanstack/react-query";

export function useRepliesByCommentID(
  isAuthenticated: boolean,
  commentId: string,
  showReplies: boolean
) {
  return useQuery({
    queryKey: ["replies", commentId],
    queryFn: () => getReplies(commentId, isAuthenticated),
    staleTime: Infinity,
    refetchInterval: Infinity,
    enabled: showReplies,
  });
}
