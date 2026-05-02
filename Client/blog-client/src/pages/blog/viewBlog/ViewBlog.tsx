import Container from "@mui/material/Container";
import Link from "@mui/material/Link";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import Typography from "@mui/material/Typography";
import { useParams } from "react-router-dom";
import BlogSection from "./BlogSection";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import SocialShareDial from "@/components/SocialShareSpeedDial/SocialShareSpeedDial";
import CommentSection from "./CommentSection";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { getBlogBySlug } from "@/api/blogApi";
import type { Blog } from "@/pages/home/BlogList";
import { CircularProgress } from "@mui/material";

export default function ViewBlog() {
  const navigate = useNavigate();
  const { slug } = useParams();

  const queryClient = useQueryClient();
  const { data: blog, isLoading } = useQuery({
    queryKey: ["blogs", slug],
    queryFn: () => getBlogBySlug(slug!),
    staleTime: Infinity,
    enabled: slug != null && slug != "",
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
        ?.find((i) => i.urlSlug == slug);
      if (!blog) return undefined;
      console.log(blog);
      return blog;
    },
  });

  if (!slug || !blog)
    return (
      <>
        <Typography>Blog not found</Typography>
        <Link component={RouterLink} to="/">
          Go back
        </Link>
      </>
    );

  const handleGoBack = () => {
    navigate(-1);
  };

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      {/* Thumbnail */}
      {!isLoading ? <BlogSection blog={blog} /> : <CircularProgress />}
      <Divider />
      {!isLoading ? (
        <Box id="comment-section">
          <CommentSection blogID={blog.blogID} slug={blog.urlSlug} />
        </Box>
      ) : (
        <CircularProgress />
      )}
      <Box
        id="footer-action"
        mt={1}
        sx={{
          display: "flex",
          justifyContent: "flex-end",
        }}
      >
        <Button onClick={handleGoBack} variant="contained">
          Go back
        </Button>
      </Box>

      <SocialShareDial />
    </Container>
  );
}
