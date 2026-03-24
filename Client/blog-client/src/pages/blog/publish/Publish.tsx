import QuoteBanner from "@/components/banner/BannerQuotes";
import Box from "@mui/material/Box";
import InputLabel from "@mui/material/InputLabel";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import { useForm, useWatch } from "react-hook-form";
import { publishBlogSchema, type PublishBlogFormValues } from "./model/schema";
import { zodResolver } from "@hookform/resolvers/zod";
import Divider from "@mui/material/Divider";
import Button from "@mui/material/Button";
import PublishIcon from "@mui/icons-material/Publish";
import CancelIcon from "@mui/icons-material/Cancel";
import { useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import List from "@mui/material/List";
import { getDirtyFieldNames } from "@/utils/mapper";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import { ClockBanner } from "@/components/banner/Clock";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { publishBlogRequest } from "@/api/blogApi";
import slugify from "slugify";

function getFieldName(fieldName: string) {
  switch (fieldName) {
    case "title":
      return "Title";
    default:
      return "Content";
  }
}

export default function PublishPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: { errors, isSubmitting, isDirty, dirtyFields },
  } = useForm<PublishBlogFormValues>({
    resolver: zodResolver(publishBlogSchema),
    defaultValues: {
      title: "",
      urlSlug: "",
      content: "",
    },
    mode: "all",
  });

  const title = useWatch({ control, name: "title" });

  useEffect(() => {
    if (title) {
      setValue(
        "urlSlug",
        slugify(title, {
          lower: true,
          strict: true,
          trim: true,
        })
      );
    }
  }, [title, setValue]);

  const { mutate, isPending } = useMutation({
    mutationFn: publishBlogRequest,
    retry: false,
    onSuccess: (data) => {
      console.log(data);
      queryClient.invalidateQueries({ queryKey: ["blogs"] });
    },
    onError: (error) => {
      if (error.message.includes("500")) {
        alert("blog url is already existed");
      }
    },
  });

  const onSubmit = async (data: PublishBlogFormValues) => {
    console.log("Form Data:", data);
    mutate(data);
  };

  const handleCancel = () => {
    if (isDirty) {
      setOpen(true);
    } else {
      navigate("/dashboard");
    }
  };

  const dirtyFieldNames = getDirtyFieldNames(dirtyFields);

  const handleConfirm = () => {};

  return (
    <>
      <QuoteBanner />
      <Box
        component="form"
        sx={{
          p: 2,
        }}
        onSubmit={handleSubmit(onSubmit)}
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
              focused
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
        <Box id="blog-content-section">
          <Box sx={{ width: "100%", p: 1 }}>
            <InputLabel htmlFor="blog-content" sx={{ mb: 1 }}>
              Content
            </InputLabel>
            <TextField
              id="blog-content"
              placeholder="What are you going to write?"
              {...register("content")}
              multiline
              fullWidth
              size="small"
              rows={4}
              error={!!errors.content}
              helperText={errors.content?.message || " "}
            />
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
            disabled={isSubmitting || isPending || !isDirty}
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
