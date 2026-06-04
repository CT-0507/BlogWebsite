import { getAuthorDashboardMetrics } from "@/api/authorApi";
import { getBlogMetrics } from "@/api/blogApi";
import { useQuery } from "@tanstack/react-query";

export function useGetAuthorDashboardData() {
  return useQuery({
    queryKey: ["dashboard", "author"],
    queryFn: getAuthorDashboardMetrics,
    staleTime: 1000 * 3600 * 2,
    refetchInterval: 1000 * 3600 * 2,
    retry: false,
  });
}

interface GetBlogMetrics {
  resultLength?: number;
  viewType: string;
  blogID: number;
  enabled: boolean;
}

export function useGetAuthorDashboardBlogMetrics({
  blogID,
  viewType,
  resultLength = 4,
  enabled,
}: GetBlogMetrics) {
  return useQuery({
    queryKey: [
      "dashboard",
      "author",
      blogID,
      "metrics",
      viewType,
      resultLength,
    ],
    queryFn: () => getBlogMetrics({ blogID, viewType, resultLength }),
    staleTime: 1000 * 3600 * 2,
    refetchInterval: 1000 * 3600 * 2,
    retry: false,
    enabled: enabled,
  });
}
