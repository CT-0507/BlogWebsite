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
import { CircularProgress } from "@mui/material";
import { useBlogBySlug } from "@/hooks/useBlogBySlug";
import BlogVoteSection from "./components/BlogVoteSection";
import { useAuth } from "@/hooks/useAuth";

export default function ViewBlog() {
  const navigate = useNavigate();
  const { slug } = useParams();
  const { isAuthenticated } = useAuth();

  const { data: blog, isLoading } = useBlogBySlug(isAuthenticated, slug);

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

      <Box id="reaction-section" sx={{ mt: 5 }}>
        {!isLoading ? (
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
        {!isLoading ? (
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
