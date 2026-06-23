import { useQueryLikedBlogs } from "@/hooks/useLikedBlogs";
import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import { Link as RouterLink } from "react-router-dom";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";
import React from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Link from "@mui/material/Link";
import CardMedia from "@mui/material/CardMedia";
import LazyImage from "@/components/Image/LazyImage";
import Tooltip from "@mui/material/Tooltip";
import placeholder from "@/assets/160x120.svg";
import { truncate } from "@/utils/textUtils";

export default function LikedBlogsPage() {
  const { data, isLoading } = useQueryLikedBlogs();

  const handleDelete = () => {
    alert("Feature is not implement yet");
  };
  return (
    <Box
      sx={{
        flex: 1,
        py: 1,
      }}
    >
      <Typography variant="h4" align="center">
        Liked Blogs
      </Typography>
      <Box sx={{ p: 1, px: 2 }}>
        <Divider />
        <List
          sx={{
            width: "100%",
            bgcolor: "background.paper",
          }}
        >
          {!isLoading &&
            (data?.total === 0 ? (
              <ListItem>Currently, you don't have any liked blogs.</ListItem>
            ) : (
              data?.blogs.map((blog, index) => (
                <React.Fragment key={index}>
                  <ListItem
                    secondaryAction={
                      <IconButton
                        edge="end"
                        aria-label="delete"
                        onClick={handleDelete}
                      >
                        <DeleteIcon />
                      </IconButton>
                    }
                  >
                    <Box width="100%">
                      <Card key={blog.blogID} sx={{ display: "flex" }}>
                        {/* Thumbnail */}
                        <Link
                          component={RouterLink}
                          to={`/blogs/${blog.urlSlug}`}
                          underline="none"
                        >
                          <CardMedia>
                            <LazyImage
                              sx={{ width: "160px", minWidth: "160px" }}
                              src={blog.thumbnailUrl ?? placeholder}
                              alt="thumbnail"
                            />
                          </CardMedia>
                        </Link>

                        {/* Content */}
                        <Box
                          sx={{
                            display: "flex",
                            flexDirection: "column",
                            flex: 1,
                          }}
                        >
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
                                component={RouterLink}
                                to={
                                  "/blogs/author/" +
                                  (blog.author.slug
                                    ? blog.author.slug
                                    : `id/${blog.author.authorID}`)
                                }
                                underline="hover"
                              >
                                {blog.author.displayName}
                              </Link>{" "}
                              •{" "}
                            </Typography>
                          </CardContent>
                        </Box>
                      </Card>
                    </Box>
                  </ListItem>
                  {index !== data?.blogs.length && <Divider />}
                </React.Fragment>
              ))
            ))}
        </List>
        <Divider />
      </Box>
    </Box>
  );
}
