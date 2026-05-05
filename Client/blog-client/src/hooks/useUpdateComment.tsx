import {
  deleteComment,
  hideComment,
  updateCommentContent,
  type GetRootCommentsResponse,
} from "@/api/blogApi";
import type { BlogComment } from "@/types/Blog";
import { useMutation, useQueryClient } from "@tanstack/react-query";

/**
 *
 * @param blogID null if comment is not root comment
 * @returns
 */
export function useUpdateCommentContent(blogID?: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateCommentContent,
    onSuccess: (data) => {
      if (blogID) {
        queryClient.setQueryData(
          ["comments", blogID],
          (old: GetRootCommentsResponse) => ({
            ...old,
            comments: old.comments.map((comment) =>
              comment.commentId === data.commentId
                ? ({
                    ...comment,
                    content: data.content,
                    updatedAt: Date.now().toString(),
                  } as BlogComment)
                : comment
            ),
          })
        );
      } else {
        queryClient.setQueryData(
          ["replies", data.commentId],
          (old: BlogComment[] = []) =>
            old.map((comment) =>
              comment.commentId === data.commentId
                ? ({
                    ...comment,
                    content: data.content,
                    updatedAt: Date.now().toString(),
                  } as BlogComment)
                : comment
            )
        );
      }
    },
  });
}

export function useHideComment(blogID?: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: hideComment,
    onSuccess: (data) => {
      if (blogID) {
        queryClient.setQueryData(
          ["comments", blogID],
          (old: GetRootCommentsResponse) => ({
            total: old.total - 1,
            comments: old.comments.map((comment) =>
              comment.commentId !== data.commentId
                ? {
                    ...comment,
                    status: "hidden",
                  }
                : comment
            ),
          })
        );
      } else {
        queryClient.setQueryData(
          ["replies", data.commentId],
          (old: BlogComment[] = []) =>
            old.map((comment) =>
              comment.commentId !== data.commentId
                ? {
                    ...comment,
                    status: "hidden",
                  }
                : comment
            )
        );
      }
    },
  });
}

export function useDeleteComment(blogID?: number) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteComment,
    onSuccess: (data) => {
      if (blogID) {
        queryClient.setQueryData(
          ["comments", blogID],
          (old: GetRootCommentsResponse) => ({
            total: old.total - 1,
            comments: old.comments.filter(
              (comment) => comment.commentId !== data.commentId
            ),
          })
        );
      } else {
        queryClient.setQueryData(
          ["replies", data.commentId],
          (old: BlogComment[] = []) =>
            old.filter((comment) => comment.commentId !== data.commentId)
        );
      }
    },
  });
}
