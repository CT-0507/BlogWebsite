import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import { postCommentSchema, type PostCommentFormValues } from "../model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { postComment } from "@/api/blogApi";
import type { BlogComment } from "../CommentSection";
import InputLabel from "@mui/material/InputLabel";
import Typography from "@mui/material/Typography";
import { useAuth } from "@/hooks/useAuth";
import MuiLink from "@mui/material/Link";
import { Link as RouterLink, useLocation, useNavigate } from "react-router-dom";

interface NewCommentProps {
  blogID: number;
  slug: string;
}

export default function NewComment({ blogID, slug }: NewCommentProps) {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const queryClient = useQueryClient();

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
      parentCommentId: undefined,
      blogID: blogID,
      depth: 0,
    },
    mode: "all",
  });
  const { mutate, isPending } = useMutation({
    mutationFn: postComment,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      queryClient.setQueryData(
        ["blogs", slug, "comments"],
        (old: BlogComment[] = []) => {
          return [...old, data];
        }
      );
    },
    onError: (error) => {
      console.log(error);
    },
  });

  const onSubmit = async (data: PostCommentFormValues) => {
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
    <Paper
      variant="outlined"
      sx={{
        p: 2,
        borderRadius: 3,
        mb: 3,
      }}
      component="form"
      onSubmit={handleSubmit(onSubmit)}
    >
      <InputLabel
        htmlFor="comment-content"
        sx={{ mb: 1, mt: 0, fontWeight: "700", opacity: 1, fontSize: "1.5rem" }}
      >
        Write a Comment
      </InputLabel>

      <Stack spacing={2}>
        <TextField
          id="comment-content"
          multiline
          minRows={3}
          fullWidth
          placeholder="Share your thoughts..."
          {...register("content")}
          error={!!errors.content}
          helperText={errors.content?.message || " "}
        />

        <Box display="flex" justifyContent="flex-end">
          <Button
            variant="contained"
            color="error"
            disabled={!isDirty}
            sx={{ mr: 1 }}
            onClick={handleClear}
          >
            Clear
          </Button>
          <Button
            variant="contained"
            type="submit"
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
      </Stack>
    </Paper>
  );
}
