import { getAuthorDashboardMetrics } from "@/api/authorApi";
import { useQuery } from "@tanstack/react-query";

export function useGetAuthorDashboardData() {
  return useQuery({
    queryKey: ["dashboard", "author"],
    queryFn: getAuthorDashboardMetrics,
    staleTime: 3600 * 2,
    refetchInterval: 3600 * 2,
  });
}
