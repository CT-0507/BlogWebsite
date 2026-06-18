import { createBlogReaction } from "@/api/blogApi";
import type { Blog } from "@/types/Blog";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function useVoteBlog(slug: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createBlogReaction,
    onSuccess: (data) => {
      queryClient.setQueryData(["blog", slug], (old: Blog): Blog => {
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
          userReaction: reaction,
          likeCount: old.likeCount + likeDelta,
          dislikeCount: old.dislikeCount + dislikeDelta,
        };
      });
    },
  });
}
