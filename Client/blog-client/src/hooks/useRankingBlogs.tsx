import { getRankingBlogs } from "@/api/blogApi";
import type { SortBlogRankingByValue, SortDir } from "@/types/types";
import { useQuery } from "@tanstack/react-query";

interface QueryBlogsParams {
  sortBy?: SortBlogRankingByValue;
  sortDir?: SortDir;
  limit?: number;
}

export default function useRankingBlogs(
  formData: QueryBlogsParams,
  page: number,
  type: "allTime" | "trending",
  enabled: boolean
) {
  return useQuery({
    queryKey: ["blogs", page, formData, type],
    queryFn: () => getRankingBlogs(formData, page, type),
    staleTime: Infinity,
    refetchInterval: Infinity,
    enabled: enabled,
  });
}
