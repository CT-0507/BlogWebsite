import ViewBlog from "@/pages/blog/viewBlog/ViewBlog";
import { useParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useBlogBySlug } from "@/hooks/useBlogBySlug";

export default function ViewBlogPage() {
  const { slug } = useParams();
  const { isAuthenticated } = useAuth();

  const { data: blog, isLoading } = useBlogBySlug(isAuthenticated, slug);
  return <ViewBlog blog={blog} isLoading={isLoading} />;
}
