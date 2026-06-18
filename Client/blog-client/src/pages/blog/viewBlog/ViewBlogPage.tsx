import ViewBlog from "@/pages/blog/viewBlog/ViewBlog";
import { useParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useBlogBySlug } from "@/hooks/useBlogBySlug";
import { appName } from "@/config/const";

export default function ViewBlogPage() {
  const { slug } = useParams();
  const { isAuthenticated } = useAuth();

  const { data: blog, isLoading } = useBlogBySlug(isAuthenticated, slug);
  const title = `${appName}  | Blog ${blog?.title}`;
  return (
    <>
      <title>{title}</title>
      <meta
        name="description"
        content="Explore articles, featured content, and updates from this author. Discover insights, follow their work, and browse their latest blog posts."
      />
      <ViewBlog blog={blog} isLoading={isLoading} />
    </>
  );
}
