import {
  createCommentReaction,
  type CreateCommentReactionResponse,
  type GetRootCommentsResponse,
} from "@/api/blogApi";
import type { BlogComment } from "@/types/Blog";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useVoteComment(
  isRootcomment: boolean,
  blogID: number,
  parentId?: string | null,
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createCommentReaction,
    onSuccess: (data: CreateCommentReactionResponse) => {
      if (isRootcomment) {
        queryClient.setQueryData(
          ["comments", blogID],
          (old: GetRootCommentsResponse): GetRootCommentsResponse => {
            if (!old) return old;

            const transitions = {
              AddLike: { delta: [1, 0], reaction: "like" },
              AddDislike: { delta: [0, 1], reaction: "dislike" },
              LikeToDislike: { delta: [-1, 1], reaction: "dislike" },
              DislikeToLike: { delta: [1, -1], reaction: "like" },
            } as const;

            const key = data.transitionType as keyof typeof transitions;

            const {
              delta: [likeDelta, dislikeDelta],
              reaction,
            } = transitions[key] ?? { delta: [0, 0], reaction: null };

            return {
              ...old,
              comments: old.comments.map((comment) =>
                comment.commentId === data.commentId
                  ? ({
                      ...comment,
                      userReaction: reaction,
                      likeCount: comment.likeCount + likeDelta,
                      dislikeCount: comment.dislikeCount + dislikeDelta,
                    } as BlogComment)
                  : comment,
              ),
            };
          },
        );
        return;
      }
      queryClient.setQueryData(
        ["replies", parentId],
        (old: BlogComment[]): BlogComment[] => {
          if (!old) return old;
          return old.map((comment) => {
            if (comment.commentId !== data.commentId) {
              return comment;
            }
            const transitions = {
              AddLike: { delta: [1, 0], reaction: "like" },
              AddDislike: { delta: [0, 1], reaction: "dislike" },
              LikeToDislike: { delta: [-1, 1], reaction: "dislike" },
              DislikeToLike: { delta: [1, -1], reaction: "like" },
            } as const;

            const key = data.type as keyof typeof transitions;

            const {
              delta: [likeDelta, dislikeDelta],
              reaction,
            } = transitions[key] ?? { delta: [0, 0], reaction: null };
            return {
              ...comment,
              userReaction: reaction,
              likeCount: comment.likeCount + likeDelta,
              dislikeCount: comment.dislikeCount + dislikeDelta,
            };
          });
        },
      );
    },
    onError: (error) => {
      console.log(error);
    },
  });
}
