import { useState } from "react";
import Avatar from "@mui/material/Avatar";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Collapse from "@mui/material/Collapse";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import ThumbUpAltOutlinedIcon from "@mui/icons-material/ThumbUpAltOutlined";
import ThumbDownAltOutlinedIcon from "@mui/icons-material/ThumbDownAltOutlined";
import ThumbDownAltIcon from "@mui/icons-material/ThumbDownAlt";
import ThumbUpAltIcon from "@mui/icons-material/ThumbUpAlt";
import ReplyOutlinedIcon from "@mui/icons-material/ReplyOutlined";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import CloseIcon from "@mui/icons-material/Close";
import { relativeTime } from "@/utils/timeUtils";
import CircularProgress from "@mui/material/CircularProgress";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import IconButton from "@mui/material/IconButton";
import Snackbar from "@mui/material/Snackbar";
import { useAuth } from "@/hooks/useAuth";
import {
  postCommentSchema,
  updateCommentContentSchema,
  type PostCommentFormValues,
  type UpdateCommentContentFormValues,
} from "../model/schema";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import MuiLink from "@mui/material/Link";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";
import { usePostComment } from "@/hooks/usePostComment";
import type { BlogComment, CommentReactionType } from "@/types/Blog";
import { useVoteComment } from "@/hooks/useVoteComment";
import { useRepliesByCommentID } from "@/hooks/useRepliesByCommentID";
import {
  useDeleteComment,
  useHideComment,
  useUpdateCommentContent,
} from "@/hooks/useUpdateComment";

interface CommentItemProps {
  comment: BlogComment;
  level: number;
  slug: string;
}

export default function CommentItem({
  comment,
  level = 0,
  slug,
}: CommentItemProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const [showReplies, setShowReplies] = useState(false);
  const [showReplyInput, setShowReplyInput] = useState(false);
  const [showError, setShowError] = useState(false);
  const [newReply, setNewReply] = useState(0);
  const [editMode, setEditMode] = useState(false);
  const [showConfirmDialog, setShowConfirmDialog] = useState(false);
  const [action, setAction] = useState<string>("");
  const [showSnackbar, setShowSnackbar] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState("");

  const { user, authLoading, isAuthenticated } = useAuth();

  const hasReplies = comment.replyCount > 0 || newReply != 0;
  const canReply = level < 2; // max reply depth

  const { data: replies, isLoading } = useRepliesByCommentID(
    isAuthenticated,
    comment.commentId,
    showReplies,
  );

  const handleShowReplies = () => {
    setShowReplies(!showReplies);
  };

  const {
    register,
    handleSubmit,
    resetField,
    formState: { errors, isSubmitting, isDirty, isValid },
  } = useForm<PostCommentFormValues>({
    resolver: zodResolver(postCommentSchema),
    defaultValues: {
      actorType: "user",
      content: "",
      parentCommentId: comment.commentId,
      rootCommentId: comment.rootCommentId,
      blogID: comment.blogId,
      depth: level + 1,
    },
    mode: "all",
  });

  const { mutate, isPending } = usePostComment();

  const onSubmit = (data: PostCommentFormValues) => {
    if (!isAuthenticated) {
      navigate("/account");
      return;
    }
    console.log("Form Data:", data);
    mutate(data, {
      onSuccess: () => {
        setNewReply((prev) => prev + 1);
        setShowReplies(true);
        resetField("content");
        setShowReplyInput(false);
      },
    });
  };

  const handleClear = () => {
    resetField("content");
  };

  const { mutate: mutateVote, isPending: isPendingVoteComment } =
    useVoteComment(level === 0, comment.blogId, comment.parentCommentId);

  const handleVote = (next: CommentReactionType) => {
    if (!isAuthenticated) {
      setShowError(true);
      return;
    }
    if (comment.userReaction === next || isPendingVoteComment) {
      return;
    }
    mutateVote({
      commentId: comment.commentId,
      userId: user!.userID,
      type: next,
    });
  };

  const handleShowEditMode = () => {
    setEditMode((prev) => !prev);
  };

  // Update comment content
  const {
    register: registerContent,
    handleSubmit: handleSubmitContent,
    resetField: resetFieldContent,
    formState: {
      errors: errorsContent,
      isSubmitting: isSubmittingContent,
      isDirty: isDirtyContent,
      isValid: isValidContent,
    },
  } = useForm<UpdateCommentContentFormValues>({
    resolver: zodResolver(updateCommentContentSchema),
    defaultValues: {
      content: comment.content,
      commentId: comment.commentId,
    },
    mode: "all",
  });

  const handleClearContent = () => {
    resetFieldContent("content");
  };

  const handleShowConfirmDialog = (action: string) => {
    setAction(action);
    setShowConfirmDialog(true);
  };

  const handleCloseConfirmDialog = () => {
    setShowConfirmDialog(false);
  };

  const { mutate: mutateHide } = useHideComment(
    level === 0 ? comment.blogId : undefined,
  );
  const { mutate: mutateDelete } = useDeleteComment(
    level === 0 ? comment.blogId : undefined,
  );

  const handleExecuteAction = () => {
    if (action === "hide") {
      mutateHide(comment.commentId, {
        onSuccess: () => {
          setShowConfirmDialog(false);
        },
        onError: () => {
          setShowSnackbar(true);
          setSnackbarMessage("Failed to hide comment. Please try again later.");
        },
      });
    } else {
      mutateDelete(comment.commentId, {
        onSuccess: () => {
          setShowConfirmDialog(false);
        },
        onError: () => {
          setShowSnackbar(true);
          setSnackbarMessage(
            "Failed to delete comment. Please try again later.",
          );
        },
      });
    }
  };

  const { mutate: mutateContent, isPending: isPendingContent } =
    useUpdateCommentContent(level === 0 ? comment.blogId : undefined);

  const onSubmitContent = (data: UpdateCommentContentFormValues) => {
    if (!isAuthenticated) {
      navigate("/account");
      return;
    }
    console.log("Form Data:", data);
    mutateContent(data, {
      onSuccess: () => {
        setEditMode(false);
      },
    });
  };

  return (
    <Box sx={{ ml: level * 4, mt: 2 }}>
      <Snackbar
        open={showSnackbar}
        autoHideDuration={6000}
        onClose={() => setShowSnackbar(false)}
        message={snackbarMessage}
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
      <Paper
        variant="outlined"
        sx={{
          p: 2,
          borderRadius: 3,
          bgcolor: "#fff",
        }}
      >
        <Stack direction="row" spacing={2}>
          <Avatar>{comment.actorAvatarUrl}</Avatar>

          <Box flex={1}>
            <Stack
              direction="row"
              spacing={1}
              alignItems="center"
              flexWrap="wrap"
            >
              <Typography fontWeight={600}>
                {comment.actorDisplayName}
              </Typography>

              <Typography variant="caption" color="text.secondary">
                {relativeTime(comment.createdAt)} ago
              </Typography>

              {comment.createdAt !== comment.updatedAt && (
                <Typography
                  variant="caption"
                  color="text.secondary"
                  sx={{ mr: 1 }}
                >
                  edited
                </Typography>
              )}

              {comment.status && comment.status === "hidden" && (
                <Typography variant="caption" color="text.secondary">
                  hidden
                </Typography>
              )}
            </Stack>

            {editMode ? (
              <Box
                component="form"
                onSubmit={handleSubmitContent(onSubmitContent)}
              >
                <Dialog
                  open={showConfirmDialog}
                  onClose={handleCloseConfirmDialog}
                >
                  <DialogTitle>Confirm</DialogTitle>
                  <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                      This action will {action} this comment. Are you sure?
                    </DialogContentText>
                  </DialogContent>
                  <DialogActions>
                    <Button onClick={handleCloseConfirmDialog} autoFocus>
                      Disagree
                    </Button>
                    <Button
                      onClick={handleExecuteAction}
                      variant="outlined"
                      color="error"
                    >
                      Agree
                    </Button>
                  </DialogActions>
                </Dialog>
                <Stack direction="column">
                  <Stack direction="row" my={1}>
                    <Button onClick={() => handleShowConfirmDialog("hide")}>
                      Hide
                    </Button>
                    <Button
                      onClick={() => handleShowConfirmDialog("delete")}
                      color="error"
                    >
                      Delete
                    </Button>
                  </Stack>
                  <TextField
                    fullWidth
                    size="small"
                    placeholder="Write a reply..."
                    {...registerContent("content")}
                    error={!!errorsContent.content}
                    helperText={errorsContent.content?.message || " "}
                  />
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "flex-start",
                    }}
                  >
                    <Button
                      variant="contained"
                      color="error"
                      disabled={!isDirty}
                      sx={{ mr: 1, flexShrink: 0 }}
                      onClick={handleClearContent}
                    >
                      Clear
                    </Button>
                    <Button
                      variant="contained"
                      type="submit"
                      sx={{
                        flexShrink: 0,
                      }}
                      disabled={
                        isPendingContent ||
                        isSubmittingContent ||
                        !isDirtyContent ||
                        !isValidContent ||
                        !isAuthenticated
                      }
                    >
                      {isPending || isSubmitting ? "Sending" : "Edit Comment"}
                    </Button>
                  </Box>
                </Stack>
              </Box>
            ) : (
              <Typography sx={{ mt: 1 }}>{comment.content}</Typography>
            )}

            <Stack
              direction="row"
              spacing={1}
              alignItems="center"
              sx={{ mt: 1 }}
              flexWrap="wrap"
            >
              <Button
                size="small"
                startIcon={
                  comment.userReaction && comment.userReaction === "like" ? (
                    <ThumbUpAltIcon />
                  ) : (
                    <ThumbUpAltOutlinedIcon />
                  )
                }
                onClick={() => handleVote("like")}
              >
                {comment.likeCount || 0}
              </Button>

              <Button
                size="small"
                startIcon={
                  comment.userReaction && comment.userReaction === "dislike" ? (
                    <ThumbDownAltIcon />
                  ) : (
                    <ThumbDownAltOutlinedIcon />
                  )
                }
                onClick={() => handleVote("dislike")}
              >
                {comment.dislikeCount || 0}
              </Button>

              {canReply && (
                <Button
                  size="small"
                  startIcon={<ReplyOutlinedIcon />}
                  onClick={() => setShowReplyInput(!showReplyInput)}
                >
                  Reply
                </Button>
              )}

              {hasReplies && (
                <Button
                  size="small"
                  endIcon={
                    showReplies ? <ExpandLessIcon /> : <ExpandMoreIcon />
                  }
                  onClick={handleShowReplies}
                >
                  {showReplies
                    ? "Hide Replies"
                    : `View Replies (${comment.replyCount})`}
                </Button>
              )}
            </Stack>

            {showError && (
              <Stack
                direction="row"
                spacing={1}
                alignItems="center"
                sx={{ mt: 1 }}
                flexWrap="wrap"
              >
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
              </Stack>
            )}

            {/* Reply Input */}
            <Collapse in={showReplyInput}>
              <Stack
                direction="row"
                spacing={1}
                sx={{ mt: 2 }}
                component="form"
                onSubmit={handleSubmit(onSubmit)}
              >
                <TextField
                  fullWidth
                  size="small"
                  placeholder="Write a reply..."
                  {...register("content")}
                  error={!!errors.content}
                  helperText={errors.content?.message || " "}
                />
                <Box
                  sx={{
                    display: "flex",
                    alignItems: "flex-start",
                  }}
                >
                  <Button
                    variant="contained"
                    color="error"
                    disabled={!isDirty}
                    sx={{ mr: 1, flexShrink: 0 }}
                    onClick={handleClear}
                  >
                    Clear
                  </Button>
                  <Button
                    variant="contained"
                    type="submit"
                    sx={{
                      flexShrink: 0,
                    }}
                    disabled={
                      isPending ||
                      isSubmitting ||
                      !isDirty ||
                      !isValid ||
                      !isAuthenticated
                    }
                  >
                    {isPending || isSubmitting ? "Sending" : "Post Comment"}
                  </Button>
                </Box>
              </Stack>
              {!isAuthenticated && (
                <Box>
                  <Typography>
                    You need to login to post comment.
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
            </Collapse>

            {/* Replies */}
            {hasReplies && (
              <Collapse in={showReplies}>
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  <Box sx={{ mt: 2 }}>
                    {replies?.map((reply) => (
                      <CommentItem
                        key={reply.commentId}
                        comment={reply}
                        slug={slug}
                        level={level + 1}
                      />
                    ))}
                  </Box>
                )}
              </Collapse>
            )}
          </Box>

          {!authLoading &&
            isAuthenticated &&
            comment.actorId === user?.userID && (
              <Box>
                <Button onClick={handleShowEditMode}>Edit</Button>
              </Box>
            )}
        </Stack>
      </Paper>
    </Box>
  );
}
