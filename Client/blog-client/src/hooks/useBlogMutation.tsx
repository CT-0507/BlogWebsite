import { publishBlogRequest, updateBlogRequest } from "@/api/blogApi";
import type { Blog } from "@/types/Blog";
import { useMutation, useQueryClient } from "@tanstack/react-query";
export const RECOVERY_KEY = "blog-recovery";

export function useBlogMutation(mode: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: mode === "create" ? publishBlogRequest : updateBlogRequest,
    retry: false,
    onSuccess: (data) => {
      queryClient.setQueryData(
        ["author_blogs", data.author.slug],
        (old: Blog[]) => [...old, data],
      );
      localStorage.remove(RECOVERY_KEY);
    },
  });
}
