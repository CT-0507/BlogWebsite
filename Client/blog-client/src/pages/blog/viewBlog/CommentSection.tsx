import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import CommentItem from "./components/CommentItem";
import NewComment from "./components/NewComment";
import Paper from "@mui/material/Paper";
import CircularProgress from "@mui/material/CircularProgress";
import { useCommentByBlogID } from "@/hooks/useCommentByBlogID";
import { useAuth } from "@/hooks/useAuth";

interface CommentSectionProps {
  blogID: number;
  slug: string;
}

export default function CommentSection({ blogID, slug }: CommentSectionProps) {
  const { isAuthenticated } = useAuth();
  const { data: comments, isLoading } = useCommentByBlogID(
    blogID,
    isAuthenticated,
  );

  return (
    <Box sx={{ mt: 4 }}>
      <Typography variant="h5" mb={2}>
        Comments ({comments?.total || 0})
      </Typography>

      <Divider sx={{ mb: 2 }} />

      <NewComment blogID={blogID} />

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
          comments?.comments!.map((comment) => (
            <CommentItem
              key={comment.commentId}
              level={0}
              slug={slug}
              comment={comment}
            />
          ))
        ) : (
          <CircularProgress />
        )}
      </Stack>
    </Box>
  );
}
