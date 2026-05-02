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
import ReplyOutlinedIcon from "@mui/icons-material/ReplyOutlined";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import ExpandLessIcon from "@mui/icons-material/ExpandLess";
import { type BlogComment } from "../CommentSection";
import { relativeTime } from "@/utils/timeUtils";
import { getReplies, postComment } from "@/api/blogApi";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { CircularProgress } from "@mui/material";
import { useAuth } from "@/hooks/useAuth";
import { postCommentSchema, type PostCommentFormValues } from "../model/schema";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import MuiLink from "@mui/material/Link";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";

interface CommentItemProps {
  comment: BlogComment;
  level: number;
  slug: string;
  repliesQueryKey: string[];
  onLike: () => void;
  onDislike: () => void;
  onReply: () => void;
}

export default function CommentItem({
  comment,
  level = 0,
  slug,
  repliesQueryKey,
  onLike,
  onDislike,
  onReply,
}: CommentItemProps) {
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();
  const [showReplies, setShowReplies] = useState(false);
  const [showReplyInput, setShowReplyInput] = useState(false);

  const { user, authLoading, isAuthenticated } = useAuth();

  const hasReplies = comment.replyCount > 0;
  const canReply = level < 2; // max reply depth

  const setQueryDataKey = [...repliesQueryKey, comment.commentId, "replies"];

  const { data: replies, isLoading } = useQuery({
    queryKey: setQueryDataKey,
    queryFn: () => getReplies(comment.commentId),
    staleTime: Infinity,
    refetchInterval: Infinity,
    enabled: showReplies,
  });

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

  const { mutate, isPending } = useMutation({
    mutationFn: postComment,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      queryClient.setQueryData(setQueryDataKey, (old: BlogComment[] = []) => {
        return [...old, data];
      });
    },
    onError: (error) => {
      console.log(error);
    },
  });

  const onSubmit = (data: PostCommentFormValues) => {
    if (!isAuthenticated) {
      navigate("/account");
      return;
    }
    console.log("Form Data:", data);
    mutate(data);
  };

  const handleClear = () => {
    resetField("content");
  };

  return (
    <Box sx={{ ml: level * 4, mt: 2 }}>
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
                <Typography variant="caption" color="text.secondary">
                  edited
                </Typography>
              )}
            </Stack>

            <Typography sx={{ mt: 1 }}>{comment.content}</Typography>

            <Stack
              direction="row"
              spacing={1}
              alignItems="center"
              sx={{ mt: 1 }}
              flexWrap="wrap"
            >
              <Button
                size="small"
                startIcon={<ThumbUpAltOutlinedIcon />}
                onClick={() => onLike()}
              >
                {comment.likes || 0}
              </Button>

              <Button
                size="small"
                startIcon={<ThumbDownAltOutlinedIcon />}
                onClick={() => onDislike()}
              >
                {comment.dislikes || 0}
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
                        repliesQueryKey={setQueryDataKey}
                        slug={slug}
                        level={level + 1}
                        onLike={onLike}
                        onDislike={onDislike}
                        onReply={onReply}
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
                <Button>Edit</Button>
              </Box>
            )}
        </Stack>
      </Paper>
    </Box>
  );
}
