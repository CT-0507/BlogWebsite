import Container from "@mui/material/Container";
import { useNavigate } from "react-router-dom";
import BlogSection from "./BlogSection";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import SocialShareDial from "@/components/SocialShareSpeedDial/SocialShareSpeedDial";
import CommentSection from "./CommentSection";
import CircularProgress from "@mui/material/CircularProgress";
import BlogVoteSection from "./components/BlogVoteSection";
import type { Blog } from "@/types/Blog";

interface ViewBlogProps {
  blog?: Blog;
  isLoading: boolean;
}

export default function ViewBlog({ blog, isLoading }: ViewBlogProps) {
  const navigate = useNavigate();

  const handleGoBack = () => {
    navigate(-1);
  };

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      {/* Thumbnail */}
      {!isLoading && blog ? <BlogSection blog={blog} /> : <CircularProgress />}
      <Divider />

      <Box id="reaction-section" sx={{ mt: 5 }}>
        {!isLoading && blog ? (
          <BlogVoteSection
            blogID={blog.blogID}
            slug={blog.urlSlug}
            likeCount={blog.likeCount || 0}
            dislikeCount={blog.dislikeCount || 0}
            userReaction={blog.userReaction}
          />
        ) : (
          <CircularProgress />
        )}
      </Box>
      <Divider />

      <Box id="comment-section" sx={{ mt: 5 }}>
        {!isLoading && blog ? (
          <CommentSection blogID={blog.blogID} slug={blog.urlSlug} />
        ) : (
          <CircularProgress />
        )}
      </Box>

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
