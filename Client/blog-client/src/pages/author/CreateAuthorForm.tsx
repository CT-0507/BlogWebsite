import Box from "@mui/material/Box";
import InputLabel from "@mui/material/InputLabel";
import TextField from "@mui/material/TextField";
import Button from "@mui/material/Button";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import ErrorIcon from "@mui/icons-material/Error";
import CloseIcon from "@mui/icons-material/Close";
import {
  createAuthorSchema,
  type CreateAuthorFormValues,
} from "./model/schema";
import { useForm, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { createAuthorRequest } from "@/api/authorApi";
import { styled } from "@mui/material/styles";
import Avatar from "@mui/material/Avatar";
import { useEffect, useRef, useState } from "react";
import { Divider, IconButton, Snackbar, Typography } from "@mui/material";
import slugify from "slugify";

const VisuallyHiddenInput = styled("input")({
  clip: "rect(0 0 0 0)",
  clipPath: "inset(50%)",
  height: 1,
  overflow: "hidden",
  position: "absolute",
  bottom: 0,
  left: 0,
  whiteSpace: "nowrap",
  width: 1,
});

export default function CreateAuthorForm() {
  const inputRef = useRef<HTMLInputElement>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarContent, setSnackbarContent] = useState("");
  const [snackbarType, setSnackbarType] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    control,
    formState: { errors, isSubmitting, isValid },
  } = useForm<CreateAuthorFormValues>({
    resolver: zodResolver(createAuthorSchema),
    defaultValues: {
      displayName: "",
      bio: "",
      avatar: null,
      slug: "",
      socialLink: "",
      email: undefined,
    },
    mode: "all",
  });

  const { mutate, isPending } = useMutation({
    mutationFn: createAuthorRequest,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      setSnackbarOpen(true);
      setSnackbarContent("Successfully created profile");
      setSnackbarType(true);
    },
    onError: (error) => {
      console.log(error);
      setSnackbarOpen(true);
      setSnackbarContent("Failed created profile");
      setSnackbarType(false);
    },
  });

  const slug = useWatch({ control, name: "displayName" });

  useEffect(() => {
    if (slug) {
      setValue(
        "slug",
        slugify(slug, {
          lower: true,
          strict: true,
          trim: true,
        })
      );
    }
  }, [slug, setValue]);

  // Cleanup
  useEffect(() => {
    return () => {
      if (preview) URL.revokeObjectURL(preview);
    };
  }, [preview]);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const objectUrl = URL.createObjectURL(file);
    setValue("avatar", file);
    setPreview(objectUrl);
  };

  const handleAvatarClick = () => {
    inputRef.current?.click();
  };

  const handleClear = () => {
    if (preview) URL.revokeObjectURL(preview);
    setPreview(null);
    setValue("avatar", null);

    // reset input so same file can be reselected
    if (inputRef.current) {
      inputRef.current.value = "";
    }
  };

  const handleCloseSnackBar = () => {
    setSnackbarOpen(false);
    setSnackbarContent("");
    setSnackbarType(false);
  };

  const onSubmit = async (data: CreateAuthorFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };
  return (
    <Box component="section" id="userEmailInfo" sx={{ p: 2 }}>
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleCloseSnackBar}
        message="Status"
        action={
          <>
            <Divider
              orientation="vertical"
              flexItem
              sx={{
                width: "2px",
                borderRight: "1px solid white",
                height: "30px",
                mr: 1,
              }}
            />
            {snackbarType ? (
              <CheckCircleIcon color="success" />
            ) : (
              <ErrorIcon color="error" />
            )}
            <Typography>{snackbarContent}</Typography>
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={handleCloseSnackBar}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          </>
        }
      />
      <Box sx={{ display: "flex", justifyContent: "flex-end" }}>
        <Typography color="red">* indicates required field</Typography>
      </Box>
      <Box component="form" onSubmit={handleSubmit(onSubmit)}>
        <Box sx={{ width: "60%" }}>
          <Box sx={{ mt: 1 }}>
            <Box position="relative" display="inline">
              <Avatar
                alt="Avatar"
                src={preview ?? ""}
                sx={{ width: 124, height: 124 }}
                onClick={handleAvatarClick}
              />
              {preview && (
                <IconButton
                  size="small"
                  onClick={handleClear}
                  sx={{
                    position: "absolute",
                    top: 0, // 👈 inside
                    left: 0, // 👈 inside
                    backgroundColor: "rgba(0,0,0,0.6)",
                    color: "white",
                    width: 24,
                    height: 24,
                    "&:hover": {
                      backgroundColor: "rgba(0,0,0,0.8)",
                    },
                  }}
                >
                  <CloseIcon fontSize="small" />
                </IconButton>
              )}
            </Box>
            <Button
              sx={{
                mt: 1,
              }}
              component="label"
              role={"upload"}
              variant="contained"
              tabIndex={-1}
              startIcon={<CloudUploadIcon />}
            >
              Upload your avatar
              <VisuallyHiddenInput
                type="file"
                onChange={handleFileChange}
                ref={inputRef}
              />
            </Button>
          </Box>
          <Box>
            <InputLabel
              htmlFor="displayName"
              sx={{ mb: 1 }}
              required
              aria-required
            >
              Display name
            </InputLabel>
            <TextField
              id="displayName"
              disabled={isSubmitting || isPending}
              placeholder="Enter your first name"
              {...register("displayName")}
              size="small"
              error={!!errors.displayName}
              helperText={errors.displayName?.message || " "}
            />
          </Box>
          <Box sx={{ mt: 1 }}>
            <InputLabel htmlFor="bio" sx={{ mb: 1 }}>
              Bio
            </InputLabel>
            <TextField
              id="bio"
              placeholder="What do you want other to know??"
              {...register("bio")}
              multiline
              fullWidth
              size="small"
              rows={4}
              error={!!errors.bio}
              helperText={errors.bio?.message || " "}
            />
          </Box>
          <Box>
            <InputLabel htmlFor="slug" sx={{ mb: 1 }} required aria-required>
              Your author url is:
            </InputLabel>
            <TextField
              id="slug"
              placeholder="Your url is"
              {...register("slug")}
              size="small"
              fullWidth
              error={!!errors.slug}
              helperText={errors.slug?.message || " "}
            />
          </Box>
          <Box>
            <InputLabel htmlFor="socialLink" sx={{ mb: 1 }}>
              Social Link
            </InputLabel>
            <TextField
              id="socialLink"
              disabled={isSubmitting || isPending}
              placeholder="Enter your social link"
              {...register("socialLink")}
              size="small"
              error={!!errors.socialLink}
              helperText={errors.socialLink?.message || " "}
              sx={{ mr: 1, width: "90%" }}
            />
          </Box>
          <Box>
            <InputLabel htmlFor="email" sx={{ mb: 1 }}>
              Email
            </InputLabel>
            <TextField
              id="email"
              disabled={isSubmitting || isPending}
              placeholder="Enter your email"
              {...register("email")}
              size="small"
              error={!!errors.email}
              helperText={errors.email?.message || " "}
              sx={{ mr: 1, width: "90%" }}
            />
          </Box>
        </Box>
        <Divider />
        <Box id="basic-info-form-action" sx={{ mt: 2 }}>
          <Button
            variant="contained"
            type="submit"
            disabled={!isValid || isPending || isSubmitting}
            sx={{
              width: "40%",
              height: "50px",
              mr: 1,
            }}
          >
            {isPending ? "Saving..." : "Save"}
          </Button>
        </Box>
      </Box>
    </Box>
  );
}
