import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import List from "@mui/material/List";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItem from "@mui/material/ListItem";
import Avatar from "@mui/material/Avatar";
import ListItemText from "@mui/material/ListItemText";
import Divider from "@mui/material/Divider";
import React from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getFollowedAuthorsRequest,
  unfollowAuthorRequest,
} from "@/api/authorApi";
import { Link } from "react-router-dom";
import IconButton from "@mui/material/IconButton";
import DeleteIcon from "@mui/icons-material/Delete";

export default function FollowedAuthorsPage() {
  const queryClient = useQueryClient();
  const { data, isLoading } = useQuery({
    queryKey: ["followed_authors"],
    queryFn: () => getFollowedAuthorsRequest(),
  });

  const { mutate, isPending } = useMutation({
    mutationFn: unfollowAuthorRequest,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["followed_authors"] });
    },
  });
  const handleUnfollowAuthor = (authorID: string) => {
    mutate(authorID);
  };
  return (
    <Box
      sx={{
        flex: 1,
        py: 1,
      }}
    >
      <Typography variant="h4" align="center">
        Followed Authors
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
            (data?.length === 0 ? (
              <ListItem>Currently, you have followed no authors.</ListItem>
            ) : (
              data?.authors.map((item, index) => (
                <React.Fragment key={index}>
                  <ListItem
                    secondaryAction={
                      <IconButton
                        edge="end"
                        aria-label="delete"
                        disabled={isPending}
                        onClick={() => handleUnfollowAuthor(item.authorID)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    }
                  >
                    <Box
                      component={Link}
                      to={`/blogs/author/${item.slug}`}
                      sx={{
                        display: "flex",
                        textDecoration: "none",
                        color: "inherit",
                      }}
                    >
                      <ListItemAvatar>
                        <Avatar src={item.avatar} />
                      </ListItemAvatar>
                      <ListItemText>{item.displayName}</ListItemText>
                    </Box>
                  </ListItem>
                  {index !== data?.authors.length && <Divider />}
                </React.Fragment>
              ))
            ))}
        </List>
        <Divider />
      </Box>
    </Box>
  );
}
