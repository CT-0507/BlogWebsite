import Link from "@mui/material/Link";
import { Link as RouterLink } from "react-router-dom";
import Typography from "@mui/material/Typography";
import type { Blog } from "@/types/Blog";
import RenderArticle from "@/components/renderBlog/RenderBlog";
import LazyImage from "@/components/Image/LazyImage";
import placeholder from "@/assets/160x120.svg";

interface BlogSectionProps {
  blog: Blog;
}

export default function BlogSection({ blog }: BlogSectionProps) {
  return (
    <>
      <LazyImage
        src={blog.thumbnailUrl ?? placeholder}
        alt="thumbnail"
        sx={{
          width: "100%",
          borderRadius: 2,
          mb: 3,
        }}
      />

      {/* Title */}
      <Typography variant="h3" gutterBottom>
        {blog.title}
      </Typography>

      {/* Author + Date */}
      <Typography variant="body2" color="text.secondary" sx={{ mb: 4 }}>
        By{" "}
        <Link
          component={RouterLink}
          to={`/blogs/author/${blog.author.slug}`}
          underline="hover"
        >
          {blog.author.displayName}
        </Link>{" "}
        • {new Date(blog.createdAt).toLocaleDateString()}
      </Typography>

      {/* Content */}
      {/* <Typography
        variant="body1"
        sx={{
          lineHeight: 1.8,
          fontSize: "1.1rem",
        }}
      >
        {blog.contentText}
      </Typography> */}
      <RenderArticle content={blog.contentJson} />
    </>
  );
}
