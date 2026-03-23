import { getBlogBySlug } from "@/api/blogApi";
import type { Blog } from "@/pages/home/BlogList";
import Link from "@mui/material/Link";
import { Link as RouterLink } from "react-router-dom";
import Typography from "@mui/material/Typography";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";

interface BlogSectionProps {
  id: string;
}
export default function BlogSection({ id }: BlogSectionProps) {
  const queryClient = useQueryClient();
  const { data: blog, isLoading } = useQuery({
    queryKey: ["blog", id],
    queryFn: () => getBlogBySlug(id),
    staleTime: Infinity,
    initialData: () => {
      // When using pagination
      //   const pages = queryClient.getQueryData(["items"]);

      //   if (!pages) return undefined;

      //   for (const page of pages.pages) {
      //     const found = page.items.find((i) => i.id === itemId);
      //     if (found) return found;
      //   }

      const blog = queryClient
        .getQueryData<Blog[]>(["blogs"])
        ?.find((i) => i.blogID == id);
      if (!blog) return undefined;
      console.log(blog);
      return blog;
    },
  });

  if (!blog) {
    return (
      <>
        <Typography>Id not found</Typography>
        <Link component={RouterLink} to="/">
          Go back
        </Link>
      </>
    );
  }
  return (
    <>
      {isLoading ? (
        <CircularProgress />
      ) : (
        <>
          <Box
            component="img"
            src="https://placehold.co/800x400"
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
              to={`/authors/${blog.author.authorID}`}
              underline="hover"
            >
              {blog.author.fullName}
            </Link>{" "}
            • {new Date(blog.createdAt).toLocaleDateString()}
          </Typography>

          {/* Content */}
          <Typography
            variant="body1"
            sx={{
              lineHeight: 1.8,
              fontSize: "1.1rem",
            }}
          >
            {blog.content}
          </Typography>
        </>
      )}
    </>
  );
}
