import { getRootComments } from "@/api/blogApi";
import { useQuery } from "@tanstack/react-query";

export function useCommentByBlogID(blogID: number, isAuthenticated?: boolean) {
  return useQuery({
    queryKey: ["comments", blogID],
    queryFn: () => getRootComments(blogID, isAuthenticated),
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });
}
