import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import CardMedia from "@mui/material/CardMedia";
import Container from "@mui/material/Container";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import Link from "@mui/material/Link";
import { Link as RouterLink } from "react-router-dom";
import Tooltip from "@mui/material/Tooltip";
import { formatDate, truncate } from "@/utils/textUtils";
import { useQuery } from "@tanstack/react-query";
import { listBlogs } from "@/api/blogApi";
import CircularProgress from "@mui/material/CircularProgress";
import Button from "@mui/material/Button";
import { relativeTime } from "@/utils/timeUtils";
import { useState } from "react";

export interface Blog {
  blogID: string;
  author: {
    authorID: string;
    nickname?: string;
    fullName?: string;
    email?: string;
  };
  title: string;
  urlSlug: string;
  content: string;
  createdAt: string;
}

export default function BlogList() {
  const [showLocaleDate, setShowLocaleDate] = useState(false);
  const { data, isLoading, refetch } = useQuery({
    queryKey: ["blogs"],
    queryFn: () => listBlogs(""),
    staleTime: 1000 * 60 * 30,
    refetchInterval: 1000 * 60 * 30,
  });
  const handleRefesh = () => {
    refetch();
  };
  const handleShowCreatedAt = () => {
    setShowLocaleDate((prev) => !prev);
  };
  return (
    <Container maxWidth="md" sx={{ mt: 4 }}>
      <Stack spacing={3}>
        <Box
          id="action"
          sx={{
            display: "flex",
            justifyContent: "flex-end",
          }}
        >
          <Button onClick={handleRefesh}>Refesh Blogs</Button>
        </Box>
        {isLoading && <CircularProgress />}
        {(data as Blog[])?.map((blog) => (
          <Card key={blog.blogID} sx={{ display: "flex" }}>
            {/* Thumbnail */}
            <Link
              component={RouterLink}
              to={`/blogs/${blog.urlSlug}`}
              underline="none"
            >
              <CardMedia
                component="img"
                sx={{ width: 160 }}
                image={`https://placehold.co/160x120`}
                alt="thumbnail"
              />
            </Link>

            {/* Content */}
            <Box sx={{ display: "flex", flexDirection: "column", flex: 1 }}>
              <CardContent>
                {/* Title */}
                <Tooltip title={blog.title}>
                  <Typography variant="h5">
                    <Link
                      component={RouterLink}
                      to={`/blogs/${blog.urlSlug}`}
                      underline="hover"
                      color="inherit"
                    >
                      {truncate(blog.title, 50)}
                    </Link>
                  </Typography>
                </Tooltip>

                {/* Author + Date */}
                <Typography
                  variant="body2"
                  color="text.secondary"
                  sx={{ mb: 1 }}
                >
                  By{" "}
                  <Link
                    href={
                      "/blogs/author/" +
                      (blog.author.nickname
                        ? blog.author.nickname
                        : `id/${blog.author.authorID}`)
                    }
                    underline="hover"
                  >
                    {blog.author.fullName}
                  </Link>{" "}
                  •{" "}
                  <Typography
                    variant="body2"
                    component="span"
                    onClick={handleShowCreatedAt}
                  >
                    {/* • {formatDate(blog.createdAt)} */}
                    {showLocaleDate
                      ? formatDate(blog.createdAt)
                      : relativeTime(blog.createdAt)}
                  </Typography>
                </Typography>

                {/* Preview */}
                <Typography variant="body1">
                  {truncate(blog.content)}
                </Typography>
              </CardContent>
            </Box>
          </Card>
        ))}
      </Stack>
    </Container>
  );
}
