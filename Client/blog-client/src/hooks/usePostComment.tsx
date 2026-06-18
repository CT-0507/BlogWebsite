import { postComment, type GetRootCommentsResponse } from "@/api/blogApi";
import type { BlogComment } from "@/types/Blog";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function usePostComment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: postComment,

    onSuccess: (data) => {
      const { blogId, parentCommentId } = data;

      // update root comments count
      if (!parentCommentId) {
        queryClient.setQueryData(
          ["comments", blogId],
          (old: GetRootCommentsResponse) => ({
            comments: [...old.comments, data],
            total: old.total + 1,
          }),
        );
      }

      // replies
      if (parentCommentId) {
        queryClient.setQueryData(
          ["replies", parentCommentId],
          (old: BlogComment[] = []) => [...old, data],
        );

        // update reply count
        queryClient.setQueryData(
          ["comments", blogId],
          (old: BlogComment[] = []) =>
            old && Array.isArray(old)
              ? old.map((c) =>
                  c.commentId === parentCommentId
                    ? { ...c, replyCount: c.replyCount + 1 }
                    : c,
                )
              : old,
        );
      }
    },
  });
}
