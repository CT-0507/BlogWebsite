import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import ThumbUpAltOutlinedIcon from "@mui/icons-material/ThumbUpAltOutlined";
import ThumbDownAltOutlinedIcon from "@mui/icons-material/ThumbDownAltOutlined";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbUpAltIcon from "@mui/icons-material/ThumbUpAlt";
import { useState } from "react";
import { useVoteBlog } from "@/hooks/useVoteBlog";
import type { BlogReaction } from "@/types/Blog";
import { useAuth } from "@/hooks/useAuth";
import MuiLink from "@mui/material/Link";
import { Link as RouterLink, useLocation } from "react-router-dom";

interface BlogVoteSectionProps {
  slug: string;
  blogID: number;
  likeCount: number;
  dislikeCount: number;
  userReaction?: string | null;
}

export default function BlogVoteSection({
  slug,
  blogID,
  likeCount,
  dislikeCount,
  userReaction,
}: BlogVoteSectionProps) {
  console.log(userReaction);
  const [reactionType, setReactionType] = useState(userReaction);
  const [showError, setShowError] = useState(false);
  const [hasMadeChanges, setHasMadeChanges] = useState(false);
  const location = useLocation();

  const { mutate, isPending } = useVoteBlog(slug);
  const { user, isAuthenticated } = useAuth();

  const handleVote = (next: "like" | "dislike") => {
    if (!isAuthenticated) {
      setShowError(true);
      return;
    }
    if (reactionType === next || isPending) {
      return;
    }
    const reaction: BlogReaction = {
      blogId: blogID,
      userId: user!.userID,
      type: next,
    };
    mutate(reaction, {
      onSuccess: () => {
        setReactionType(next);
        setHasMadeChanges(true);
      },
    });
  };
  return (
    <Box sx={{ pb: 3 }}>
      <Typography sx={{ mt: 1 }} variant="h4">
        How do you rate this blog ?
      </Typography>
      <Stack
        direction="row"
        spacing={1}
        alignItems="center"
        sx={{ mt: 1 }}
        flexWrap="wrap"
      >
        <Button
          size="large"
          variant="contained"
          startIcon={
            reactionType && reactionType === "like" ? (
              <ThumbUpAltIcon />
            ) : (
              <ThumbUpAltOutlinedIcon />
            )
          }
          onClick={() => handleVote("like")}
        >
          {likeCount || 0}
        </Button>

        <Button
          size="large"
          variant="contained"
          startIcon={
            reactionType && reactionType === "dislike" ? (
              <ThumbDownAltIcon />
            ) : (
              <ThumbDownAltOutlinedIcon />
            )
          }
          onClick={() => handleVote("dislike")}
        >
          {dislikeCount || 0}
        </Button>
        {showError && (
          <Box>
            <Typography>
              You need to login to vote blog.
              <MuiLink
                component={RouterLink}
                to={`/account`}
                state={{ from: location }}
                underline="hover"
              >
                To login page.
              </MuiLink>
            </Typography>
          </Box>
        )}
        {hasMadeChanges && (
          <Box>
            <Typography>Thank you for your feed back.</Typography>
          </Box>
        )}
      </Stack>
    </Box>
  );
}
