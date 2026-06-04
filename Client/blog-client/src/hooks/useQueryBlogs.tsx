import { listBlogs, listMyBlogs } from "@/api/blogApi";
import type { SortByValue, SortDir } from "@/types/types";
import { useQuery } from "@tanstack/react-query";

interface QueryBlogsParams {
  title?: string | null;
  content?: string | null;
  author?: string | null;
  sortBy: SortByValue;
  sortDir: SortDir;
  limit: number;
}

export function useQueryBlogs(formData: QueryBlogsParams, page: number) {
  return useQuery({
    queryKey: ["blogs", page, formData],
    queryFn: () => listBlogs(formData, page),
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });
}

interface QueryBlogsParams {
  title?: string | null;
  content?: string | null;
  author?: string | null;
  sortBy: SortByValue;
  sortDir: SortDir;
  limit: number;
}

export function useQueryBlogsAuthor(
  formData: Omit<QueryBlogsParams, "author">,
  page: number,
) {
  return useQuery({
    queryKey: ["blogs", page, formData],
    queryFn: () => listMyBlogs(formData, page),
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });
}
