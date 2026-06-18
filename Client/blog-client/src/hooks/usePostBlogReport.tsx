import { createBlogReport } from "@/api/blogApi";
import { useMutation } from "@tanstack/react-query";

export function usePostBlogReport() {
  return useMutation({
    mutationFn: createBlogReport,
  });
}
