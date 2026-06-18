import { useParams } from "react-router-dom";
import BlogForm from "../components/BlogForm";
import { useBlogBySlug } from "@/hooks/useBlogBySlug";
import { useAuth } from "@/hooks/useAuth";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import { appName } from "@/config/const";

export default function EditBlog() {
  const { slug } = useParams();
  const { isAuthenticated } = useAuth();

  const title = `${appName} | Edit | {blog?.title}`;

  const { data: blog, isLoading } = useBlogBySlug(isAuthenticated, slug);

  if (isLoading)
    return (
      <Box sx={{ height: "100%", width: "100%" }}>
        <CircularProgress />
      </Box>
    );

  return (
    <>
      <title>{title}</title>
      <BlogForm blog={blog} mode="edit" />
    </>
  );
}
