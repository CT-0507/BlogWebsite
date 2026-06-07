import { getNotifications } from "@/api/userApi";
import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Divider from "@mui/material/Divider";
import Fade from "@mui/material/Fade";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Popper from "@mui/material/Popper";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useRef, useState, useMemo, useEffect, Fragment } from "react";
import NotificationsIcon from "@mui/icons-material/Notifications";
import DeleteIcon from "@mui/icons-material/Delete";
import MarkAsUnreadIcon from "@mui/icons-material/MarkAsUnread";
import ClickAwayListener from "@mui/material/ClickAwayListener";
import Stack from "@mui/material/Stack";
import Link from "@mui/material/Link";
import { Link as RouterLink } from "react-router-dom";

interface NotificationContent {
  AuthorID: string;
  AuthorName: string;
  AuthorSlug: string;
  Content: string;
  Title: string;
  UrlSlug: string;
}

interface Notification {
  notificationId: string;
  content: NotificationContent;
  isRead?: boolean;
  createdAt: string;
}

export default function NotificationMenu() {
  const baseTitleRef = useRef(document.title);
  const queryClient = useQueryClient();
  const [open, setOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [expand, setExpand] = useState(false);
  const { data, isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: getNotifications,
    retry: false,
    // refetchInterval: 30 * 60 * 1000,
  });

  const defaultShowNotificationNumber = 6;

  const nofiticationCount = useMemo(() => {
    if (!isLoading && data && Array.isArray(data)) {
      return data.filter((item) => !item.isRead).length;
    }
    return 0;
  }, [data, isLoading]);

  const handleShowNotifications = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
    setOpen((previousOpen) => !previousOpen);
  };

  const handleExpand = () => {
    setExpand((prev) => !prev);
  };

  useEffect(() => {
    const baseTitle = baseTitleRef.current;

    document.title =
      nofiticationCount !== 0
        ? `${baseTitle} (${nofiticationCount})`
        : baseTitle;
  }, [nofiticationCount]);

  const canBeOpen = open && Boolean(anchorEl);
  const id = canBeOpen ? "transition-popper" : undefined;

  const handleDelete = (notificationId: string) => {
    queryClient.setQueryData(["notifications"], (old: Notification[]) =>
      old.filter((item) => item.notificationId !== notificationId),
    );
  };
  const handleMarkAsUnRead = (notificationId: string) => {
    console.log(notificationId);
  };

  return (
    <>
      <Tooltip title="Notification">
        <IconButton
          size="large"
          aria-label={`show ${nofiticationCount} new notifications`}
          color="inherit"
          onClick={handleShowNotifications}
        >
          <Badge
            badgeContent={isLoading ? <CircularProgress /> : nofiticationCount}
            color="error"
          >
            <NotificationsIcon />
          </Badge>
        </IconButton>
      </Tooltip>
      <Popper id={id} open={open} anchorEl={anchorEl} transition>
        {({ TransitionProps }) => (
          <ClickAwayListener onClickAway={() => setOpen(false)}>
            <Fade {...TransitionProps} timeout={350}>
              <Box
                sx={{
                  border: 1,
                  py: 1,
                  borderColor: "black",
                  bgcolor: "black",
                  color: "white",
                  borderRadius: 2,
                  width: "400px",
                }}
              >
                <List
                  disablePadding
                  sx={{
                    maxHeight: "550px",
                    overflowY: "auto",

                    /* Firefox */
                    scrollbarWidth: "thin",
                    scrollbarColor: "gray transparent",

                    /* Chrome, Edge, Safari */
                    "&::-webkit-scrollbar": {
                      width: "8px",
                    },
                    "&::-webkit-scrollbar-track": {
                      background: "transparent",
                    },
                    "&::-webkit-scrollbar-thumb": {
                      backgroundColor: "gray",
                      borderRadius: "8px",
                    },
                  }}
                >
                  <ListItem disablePadding sx={{ px: 1 }}>
                    <Typography variant="h6">Notification</Typography>
                  </ListItem>
                  <Divider variant="fullWidth" sx={{ bgcolor: "gray" }} />
                  {/* <ListItem disablePadding sx={{ px: 1 }}>
                  <Typography variant="h6">Notification</Typography>
                </ListItem> */}
                  {data &&
                    Array.isArray(data) &&
                    (data as Notification[])
                      .slice(
                        0,
                        expand ? data.length : defaultShowNotificationNumber,
                      )
                      .map((item) => (
                        <Fragment key={item.notificationId}>
                          <Tooltip
                            title={`${item.content.AuthorName} has created a new blog`}
                            placement="bottom-start"
                          >
                            <ListItem
                              disablePadding
                              sx={{
                                px: 1,
                              }}
                            >
                              <Stack>
                                <Typography
                                  variant="h6"
                                  sx={{
                                    overflow: "hidden",
                                    textOverflow: "ellipsis",
                                    whiteSpace: "nowrap",
                                  }}
                                >
                                  <Link
                                    component={RouterLink}
                                    to={
                                      "/blogs/author/" + item.content.AuthorSlug
                                    }
                                    underline="hover"
                                  >
                                    {item.content.AuthorName}
                                  </Link>{" "}
                                  has created a new blog
                                </Typography>
                                <Typography
                                  sx={{
                                    overflow: "hidden",
                                    textOverflow: "ellipsis",
                                    whiteSpace: "nowrap",
                                  }}
                                >
                                  Title:{" "}
                                  <Link
                                    component={RouterLink}
                                    to={`/blogs/${item.content.UrlSlug}`}
                                  >
                                    {item.content.Title}
                                  </Link>
                                </Typography>
                                <Typography>{item.content.Content}</Typography>
                              </Stack>
                              <Tooltip title="Marked as read">
                                <IconButton
                                  color="info"
                                  onClick={() =>
                                    handleMarkAsUnRead(item.notificationId)
                                  }
                                >
                                  <MarkAsUnreadIcon />
                                </IconButton>
                              </Tooltip>
                              <Tooltip title="Delete">
                                <IconButton
                                  color="error"
                                  onClick={() =>
                                    handleDelete(item.notificationId)
                                  }
                                >
                                  <DeleteIcon />
                                </IconButton>
                              </Tooltip>
                            </ListItem>
                          </Tooltip>
                          <Divider
                            variant="fullWidth"
                            sx={{ bgcolor: "gray" }}
                          />
                        </Fragment>
                      ))}
                  <ListItem
                    disablePadding
                    sx={{ placeContent: "center", pt: 1 }}
                    onClick={handleExpand}
                  >
                    {expand ? "Show less" : "Expand"}...
                  </ListItem>
                </List>
              </Box>
            </Fade>
          </ClickAwayListener>
        )}
      </Popper>
    </>
  );
}
