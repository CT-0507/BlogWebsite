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
import FlagIcon from "@mui/icons-material/Flag";
import EmojiFlagsIcon from "@mui/icons-material/EmojiFlags";
import Dialog from "@mui/material/Dialog";
import CloseIcon from "@mui/icons-material/Close";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import InputLabel from "@mui/material/InputLabel";
import Snackbar from "@mui/material/Snackbar";
import TextField from "@mui/material/TextField";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { reportSchema, type BlogReportFormValues } from "../model/schema";
import { usePostBlogReport } from "@/hooks/usePostBlogReport";

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
  const [reactionType, setReactionType] = useState(userReaction);
  const [hasReport, setHasReport] = useState(false);
  const [showError, setShowError] = useState(false);
  const [hasMadeChanges, setHasMadeChanges] = useState(false);
  const [showReportDialog, setShowReportDialog] = useState(false);
  const [showSnackbar, setShowSnackbar] = useState(false);
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

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<BlogReportFormValues>({
    resolver: zodResolver(reportSchema),
    defaultValues: {
      blogID: blogID,
      reason: "",
    },
    mode: "all",
  });

  const handleReportDialog = () => {
    setShowReportDialog(true);
  };

  const handleCloseReportDialog = () => {
    setShowReportDialog(false);
  };

  const { mutate: mutateReport, isPending: isPendingReport } =
    usePostBlogReport();

  const onSubmitReport = async (data: BlogReportFormValues) => {
    if (!isAuthenticated) {
      setShowError(true);
      return;
    }

    mutateReport(data, {
      onSuccess: () => {
        setHasReport(true);
        setShowReportDialog(false);
      },
    });
  };
  return (
    <Box sx={{ pb: 3 }}>
      <Snackbar
        open={showSnackbar}
        autoHideDuration={6000}
        onClose={() => setShowSnackbar(false)}
        message={
          "Your report has successfully submited. Thank you for your contribution"
        }
        action={
          <>
            <Button
              color="secondary"
              size="small"
              onClick={() => setShowSnackbar(false)}
            >
              UNDO
            </Button>
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={() => setShowSnackbar(false)}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          </>
        }
      />
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

        <Button
          size="large"
          variant="contained"
          sx={{ ml: 2 }}
          disabled={isSubmitting || isPendingReport}
          startIcon={hasReport ? <FlagIcon /> : <EmojiFlagsIcon />}
          onClick={() => handleReportDialog()}
        >
          Report
        </Button>
        {showError && (
          <Box>
            <Typography>
              You need to login to vote blog or report.
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

      <Dialog open={showReportDialog} onClose={handleCloseReportDialog}>
        <DialogTitle>Create report for this blog</DialogTitle>
        <Box component="form" onSubmit={handleSubmit(onSubmitReport)}>
          <DialogContent>
            <InputLabel>State your reason</InputLabel>
            <TextField
              id="report-reason"
              multiline
              minRows={3}
              fullWidth
              placeholder="Why do you want to report this blog?"
              {...register("reason")}
              error={!!errors.reason}
              helperText={errors.reason?.message || " "}
            />
          </DialogContent>
          <DialogActions>
            <Button type="submit">Submit</Button>
            <Button color="error" onClick={handleCloseReportDialog}>
              Discard & Leave
            </Button>
          </DialogActions>
        </Box>
      </Dialog>
    </Box>
  );
}
