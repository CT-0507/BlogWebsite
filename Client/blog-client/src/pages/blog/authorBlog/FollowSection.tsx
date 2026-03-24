import {
  followAuthorRequest,
  getFollowedAuthorsRequest,
  unfollowAuthorRequest,
} from "@/api/authorApi";
import type { Author, FollowedAuthorReponse } from "@/types/types";
import Button from "@mui/material/Button";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

interface FollowSectionProps {
  author: Author;
}

export default function FollowSection({ author }: FollowSectionProps) {
  const queryClient = useQueryClient();
  const { data } = useQuery({
    queryKey: ["followed_authors"],
    queryFn: getFollowedAuthorsRequest,
  });

  const isFollowed = !!data?.authors?.includes(author?.authorID);

  const toggleFollowMutation = useMutation({
    mutationFn: async ({
      authorId,
      follow,
    }: {
      authorId: string;
      follow: boolean;
    }) => {
      return follow
        ? followAuthorRequest(authorId)
        : unfollowAuthorRequest(authorId);
    },

    onMutate: async ({ authorId, follow }) => {
      await queryClient.cancelQueries({ queryKey: ["followed_authors"] });

      const previous = queryClient.getQueryData<FollowedAuthorReponse>([
        "followed_authors",
      ]);

      queryClient.setQueryData<FollowedAuthorReponse>(
        ["followed_authors"],
        (old) => {
          if (!old) return old;

          const authors = old.authors ?? [];

          return {
            ...old,
            authors: follow
              ? [...authors, authorId] // ✅ add
              : authors.filter((id) => id !== authorId), // ✅ remove
          };
        }
      );

      queryClient.setQueryData(["author", author.slug], (old?: Author) => {
        if (!old) return old;

        return {
          ...old,
          followerCount: old.followerCount + (follow ? 1 : -1),
        };
      });

      return { previous };
    },
    onSuccess: (_data, variables) => {
      queryClient.setQueryData(
        ["author", author.slug],
        (old: Author | undefined) => {
          if (!old) return old;

          return {
            ...old,
            followerCount: old.followerCount + (variables.follow ? 1 : -1),
          };
        }
      );
    },

    onError: (_err, _vars, context) => {
      if (context?.previous) {
        queryClient.setQueryData(["followed_authors"], context.previous);
      }
    },

    onSettled: () => {
      queryClient.invalidateQueries({ queryKey: ["followed_authors"] });
    },
  });

  //   const { mutate: fMutate, isPending: fLoading } = useMutation({
  //     mutationFn: followAuthorRequest,
  //     retry: false,
  //     onSuccess: (data) => {
  //       console.log(data);
  //       queryClient.setQueryData(["author", author.slug], (old: Author[]) => {
  //         if (!old) return old;
  //         return old.map((item) =>
  //           item.slug === author.slug
  //             ? {
  //                 ...item,
  //                 ["followerCount"]: item["followerCount"] + 1,
  //               }
  //             : item
  //         );
  //       });
  //     },
  //     onError: (error) => {
  //       console.log(error);
  //     },
  //   });

  //   const { mutate: uMutate, isPending: uLoading } = useMutation({
  //     mutationFn: unfollowAuthorRequest,
  //     retry: false,
  //     onSuccess: (data) => {
  //       console.log(data);
  //       queryClient.setQueryData(["author", author.slug], (old: Author[]) => {
  //         if (!old) return old;
  //         return old.map((item) =>
  //           item.slug === author.slug
  //             ? {
  //                 ...item,
  //                 ["followerCount"]: item["followerCount"] - 1,
  //               }
  //             : item
  //         );
  //       });
  //       queryClient.setQueryData(["followed_authors"])
  //     },
  //     onError: (error) => {
  //       console.log(error);
  //     },
  //   });
  //   console.log(
  //     queryClient.getQueryData(["followed_authors"]) as FollowedAuthorReponse
  //   );
  //   const handleFollow = () => {
  //     fMutate(author.authorID);
  //   };
  //   const handleUnfollw = () => {
  //     uMutate(author.authorID);
  //   };
  return (
    <>
      {/* {isFollowed ? (
        <Button variant="outlined" onClick={handleUnfollw} disabled={uLoading}>
          Unfollow
        </Button>
      ) : (
        <Button variant="outlined" onClick={handleFollow} disabled={fLoading}>
          Follow
        </Button>
      )} */}
      <Button
        onClick={() =>
          toggleFollowMutation.mutate({
            authorId: author.authorID,
            follow: !isFollowed,
          })
        }
      >
        {isFollowed ? "Unfollow" : "Follow"}
      </Button>
    </>
  );
}
