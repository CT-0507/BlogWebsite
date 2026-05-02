import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import CommentItem from "./components/CommentItem";
import NewComment from "./components/NewComment";
import { useQuery } from "@tanstack/react-query";
import { getRootComments } from "@/api/blogApi";
import Paper from "@mui/material/Paper";
import { CircularProgress } from "@mui/material";

export interface BlogComment {
  commentId: string;
  blogId: number;
  actorId: string;
  actorType: string;
  actorAvatarUrl?: string | null;
  actorDisplayName: string;
  content: string;
  likes: number;
  dislikes: number;
  replyCount: number;
  rootCommentId: string;
  createdAt: string;
  updatedAt: string;
  status?: string | null;
  parentCommentId?: string | null;
}

interface CommentSectionProps {
  blogID: number;
  slug: string;
}

export default function CommentSection({ blogID, slug }: CommentSectionProps) {
  const { data: comments, isLoading } = useQuery({
    queryKey: ["blogs", slug, "comments"],
    queryFn: () => getRootComments(blogID),
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });

  const handleLike = () => {};

  const handleDislike = () => {};

  const handleReply = () => {};

  const repliesQueryKey = ["blogs", slug, "comments"];

  return (
    <Box sx={{ mt: 4 }}>
      <Typography variant="h5" mb={2}>
        Comments ({comments?.length || 0})
      </Typography>

      <Divider sx={{ mb: 2 }} />

      <NewComment blogID={blogID} slug={slug} />

      <Divider sx={{ mb: 2 }} />

      <Stack spacing={2}>
        {!comments && (
          <Box sx={{ ml: 0 * 4, mt: 2 }}>
            <Paper
              variant="outlined"
              sx={{
                p: 2,
                borderRadius: 3,
                bgcolor: "#fff",
              }}
            >
              <Typography>Be the first to comment</Typography>
            </Paper>
          </Box>
        )}
        {!isLoading ? (
          comments!.map((comment) => (
            <CommentItem
              key={comment.commentId}
              level={0}
              slug={slug}
              repliesQueryKey={repliesQueryKey}
              comment={comment}
              onLike={handleLike}
              onDislike={handleDislike}
              onReply={handleReply}
            />
          ))
        ) : (
          <CircularProgress />
        )}
      </Stack>
    </Box>
  );
}
