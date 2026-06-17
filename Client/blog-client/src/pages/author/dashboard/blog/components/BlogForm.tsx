import QuoteBanner from "@/components/banner/BannerQuotes";
import Box from "@mui/material/Box";
import InputLabel from "@mui/material/InputLabel";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  Controller,
  useForm,
  useWatch,
  type ControllerRenderProps,
} from "react-hook-form";
import {
  editorContentSchema,
  publishBlogSchema,
  type PublishBlogFormValues,
} from "../publish/model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import Divider from "@mui/material/Divider";
import Button from "@mui/material/Button";
import PublishIcon from "@mui/icons-material/Publish";
import CancelIcon from "@mui/icons-material/Cancel";
import { useNavigate } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import List from "@mui/material/List";
import { getDirtyFieldNames } from "@/utils/mapper";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ClockBanner from "@/components/banner/Clock";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  publishBlogRequest,
  updateBlogRequest,
  type UpdateBlogRequestParams,
} from "@/api/blogApi";
import slugify from "slugify";
import type { Blog } from "@/types/Blog";
import Stack from "@mui/material/Stack";
import Chip from "@mui/material/Chip";
import Editor, { type EditorHandle } from "./EditorField";
import type { OutputData } from "@editorjs/editorjs";
import FormControl from "@mui/material/FormControl";
import FormHelperText from "@mui/material/FormHelperText";
import axios from "axios";

function getFieldName(fieldName: string) {
  switch (fieldName) {
    case "title":
      return "Title";
    default:
      return "Content";
  }
}
interface ImageFieldProps {
  field: ControllerRenderProps<PublishBlogFormValues, "thumbnail">;
}
function ImageField({ field }: ImageFieldProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement | null>(null);
  useEffect(() => {
    if (!field.value) {
      // eslint-disable-next-line react-hooks/set-state-in-effect
      setPreviewUrl(null);
      return;
    }

    const objectUrl = URL.createObjectURL(field.value);

    setPreviewUrl(objectUrl);

    return () => {
      URL.revokeObjectURL(objectUrl);
    };
  }, [field.value]);

  return (
    <Box>
      <Typography variant="subtitle1" gutterBottom>
        Upload Image
      </Typography>

      {/* Upload Button */}
      <Button variant="outlined" component="label">
        Choose File
        <input
          ref={inputRef}
          hidden
          type="file"
          accept="image/*"
          onChange={(e) => {
            const file = e.target.files?.[0] || null;

            field.onChange(file);
            e.target.value = "";
          }}
        />
      </Button>

      {/* Preview */}
      {previewUrl && (
        <Box
          sx={{
            mt: 2,
            width: 240,
            height: 240,
            position: "relative",
            border: "1px dashed #ccc",
            borderRadius: 2,
            overflow: "hidden",
            background: "#fafafa",
          }}
        >
          {/* Remove Button */}
          <Button
            size="small"
            onClick={() => field.onChange(null)}
            sx={{
              minWidth: 0,
              width: 32,
              height: 32,
              position: "absolute",
              top: 8,
              right: 8,
              borderRadius: "50%",
              background: "rgba(0,0,0,0.6)",
              color: "#fff",
              zIndex: 1,

              "&:hover": {
                background: "rgba(0,0,0,0.8)",
              },
            }}
          >
            ✕
          </Button>

          {/* Image */}
          <Box
            component="img"
            src={previewUrl}
            alt="Preview"
            sx={{
              width: "100%",
              height: "100%",
              objectFit: "cover",
            }}
          />
        </Box>
      )}
    </Box>
  );
}

export const RECOVERY_KEY = "blog-recovery";

interface BlogFormProps {
  blog?: Blog;
  mode: "create" | "edit";
}

export default function PublishPage({ blog, mode }: BlogFormProps) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);
  const [tagInput, setTagInput] = useState("");

  const editorRef = useRef<EditorHandle>(null);
  const keyRef = useRef(crypto.randomUUID());

  const {
    register,
    handleSubmit,
    control,
    setValue,
    reset,
    formState: { errors, isSubmitting, isDirty, dirtyFields },
  } = useForm<PublishBlogFormValues>({
    resolver: zodResolver(publishBlogSchema),
    defaultValues: {
      title: blog?.title ?? "",
      urlSlug: blog?.urlSlug ?? "",
      content: {
        json: blog?.contentJson ?? { blocks: [] },
        plainText: blog?.contentText ?? "",
      },
      tags: blog?.tags || [],
      thumbnail: null,
    },
    mode: "all",
  });

  const tags = useWatch({ control, name: "tags" });

  const title = useWatch({ control, name: "title" });

  useEffect(() => {
    if (title) {
      setValue(
        "urlSlug",
        slugify(title, {
          lower: true,
          strict: true,
          trim: true,
        }),
      );
    }
  }, [title, setValue]);

  const { mutate, isPending } = useMutation({
    mutationFn: mode === "create" ? publishBlogRequest : updateBlogRequest,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      reset();
      queryClient.invalidateQueries({ queryKey: ["blogs"] });
      localStorage.removeItem(RECOVERY_KEY);
      navigate("/author/my-blogs");
    },
    onError: (error) => {
      if (axios.isAxiosError(error)) {
        const status = error.response?.status;

        // Generate new key only when server definitely responded
        // and the request is not still processing.
        if (status && status !== 409) {
          keyRef.current = crypto.randomUUID();
        }

        if (status === 500) {
          alert("Blog URL already exists");
        }
      }
    },
  });

  function extractPlainText(data?: OutputData | null): string {
    if (!data || !data?.blocks) return "";

    return data.blocks
      .map((block) => {
        switch (block.type) {
          case "paragraph":
          case "header":
            return block.data.text;

          default:
            return "";
        }
      })
      .join(" ")
      .replace(/<[^>]*>/g, "")
      .trim();
  }

  const [editorError, setEditorError] = useState("");

  const onSubmit = async (data: PublishBlogFormValues) => {
    console.log("Form Data:", data);
    let saveData: UpdateBlogRequestParams;
    try {
      const editorData = await editorRef.current?.save();
      if (
        !editorData ||
        !editorData.content ||
        editorData.content.blocks.length === 0
      ) {
        alert("Null");
        return;
      }
      setEditorError("");
      const plainText = extractPlainText(editorData.content);
      const validation = editorContentSchema.safeParse({
        plainText: plainText,
      });
      if (!validation.success) {
        setEditorError(validation.error.issues[0]?.message);

        return;
      }
      saveData = {
        ...data,
        content: {
          json: editorData.content,
          plainText: plainText,
        },
        files: editorData.files,
        blogId: blog?.blogID ?? 0,
        idempotencyKey: keyRef.current,
      };

      console.log(saveData);
      mutate(saveData, {
        onError: (error) => {
          keyRef.current = crypto.randomUUID();
          if (error.message.includes("500")) {
            alert("blog url is already existed");
          }
        },
      });
    } catch (err) {
      console.log(err);
    }
  };

  const handleCancel = () => {
    if (isDirty) {
      setOpen(true);
    } else {
      navigate("/author/dashboard");
    }
  };

  const handleAddTag = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && tagInput.trim()) {
      e.preventDefault();

      const trimmed = tagInput.trim();

      if (!tags!.includes(trimmed)) {
        setValue("tags", [...tags!, trimmed], {
          shouldValidate: true,
        });
      }

      setTagInput("");
    }
  };

  const handleDeleteTag = (tagToDelete: string) => {
    setValue(
      "tags",
      tags!.filter((tag) => tag !== tagToDelete),
      { shouldValidate: true },
    );
  };

  const dirtyFieldNames = getDirtyFieldNames(dirtyFields);

  const handleConfirm = () => {
    navigate(-1);
  };
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const onErr = (err: any) => {
    console.log(err);
  };
  return (
    <>
      <QuoteBanner />
      <Box
        component="form"
        sx={{
          p: 2,
        }}
        // eslint-disable-next-line react-hooks/refs
        onSubmit={handleSubmit(onSubmit, onErr)}
      >
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            flexDirection: "column",
            p: 1,
          }}
        >
          <Box alignSelf="flex-end">
            <ClockBanner />
          </Box>
          <Typography variant="h4">Let publish new blog</Typography>
        </Box>
        <Box>
          <Controller
            name="thumbnail"
            control={control}
            render={({ field }) => <ImageField field={field} />}
          />
        </Box>
        <Box id="blog-title-section">
          <Box sx={{ width: "45%", p: 1 }}>
            <InputLabel htmlFor="blog-title" sx={{ mb: 1 }}>
              Your blog title is:
            </InputLabel>
            <TextField
              id="blog-title"
              placeholder="Title"
              {...register("title")}
              size="small"
              fullWidth
              error={!!errors.title}
              helperText={errors.title?.message || " "}
            />
          </Box>
        </Box>
        <Box id="blog-url-section">
          <Box sx={{ width: "45%", p: 1 }}>
            <InputLabel htmlFor="blog-url" sx={{ mb: 1 }}>
              Your blog url is:
            </InputLabel>
            <TextField
              id="blog-url"
              placeholder="url"
              {...register("urlSlug")}
              size="small"
              focused
              fullWidth
              error={!!errors.urlSlug}
              helperText={errors.urlSlug?.message || " "}
            />
          </Box>
        </Box>
        {/* ----------------------------- */}
        {/* Tags */}
        {/* ----------------------------- */}

        <Box>
          <Typography variant="subtitle1" gutterBottom>
            Tags
          </Typography>

          <Stack direction="row" spacing={1} flexWrap="wrap" sx={{ mb: 1 }}>
            {tags &&
              tags.map((tag) => (
                <Chip
                  key={tag}
                  label={tag}
                  onDelete={() => handleDeleteTag(tag)}
                />
              ))}
          </Stack>

          <TextField
            fullWidth
            label="Add tag"
            value={tagInput}
            onChange={(e) => setTagInput(e.target.value)}
            onKeyDown={handleAddTag}
          />
        </Box>
        <Box id="blog-content-section" sx={{ mt: 1 }}>
          <InputLabel sx={{ mb: 1, display: "block" }}>Content</InputLabel>
          <Box sx={{ width: "100%", p: 0 }}>
            <FormControl fullWidth>
              {/* <TextField
              id="blog-content"
              placeholder="What are you going to write?"
              {...register("content")}
              multiline
              fullWidth
              size="small"
              rows={4}
              error={!!errors.content}
              helperText={errors.content?.message || " "}
            /> */}
              <Editor ref={editorRef} initialData={blog?.contentJson} />
              <FormHelperText error={editorError !== ""}>
                {editorError || " "}
              </FormHelperText>
            </FormControl>
          </Box>
        </Box>
        <Divider />
        <Box
          id="form-action"
          sx={{
            p: 1,
            pt: 2,
            display: "flex",
            justifyContent: "space-around",
          }}
        >
          <Button
            type="submit"
            variant="contained"
            sx={{
              width: "45%",
              height: "50px",
            }}
            disabled={isSubmitting || isPending}
          >
            {isPending ? (
              "Loading"
            ) : (
              <>
                <PublishIcon />
                Publish
              </>
            )}
          </Button>
          <Button
            color="error"
            variant="contained"
            onClick={handleCancel}
            sx={{
              width: "45%",
              height: "50px",
            }}
          >
            <CancelIcon />
            Cancel
          </Button>
        </Box>
        <Dialog open={open} onClose={() => setOpen(false)}>
          <DialogTitle>Discard changes?</DialogTitle>

          <DialogContent>
            <Typography sx={{ mb: 1 }}>
              You have unsaved changes in the following fields:
            </Typography>

            <List dense>
              {dirtyFieldNames.map((field) => (
                <ListItem key={field as string}>
                  <ListItemText primary={getFieldName(field as string)} />
                </ListItem>
              ))}
            </List>
          </DialogContent>

          <DialogActions>
            <Button onClick={() => setOpen(false)}>Stay</Button>

            <Button color="error" onClick={handleConfirm}>
              Discard & Leave
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </>
  );
}
