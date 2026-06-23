import { listLikedBlogs } from "@/api/blogApi";
import { useQuery } from "@tanstack/react-query";

export function useQueryLikedBlogs() {
  return useQuery({
    queryKey: ["me", "liked_blogs"],
    queryFn: listLikedBlogs,
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });
}
