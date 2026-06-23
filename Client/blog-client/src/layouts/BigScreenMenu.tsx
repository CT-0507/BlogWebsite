import { logoutRequest } from "@/api/auth";
import { tokenStore } from "@/api/store/tokenStore";
import CircularProgress from "@mui/material/CircularProgress";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import { useQueryClient, useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";

interface BigScreenMenuProps {
  menuId: string;
  anchorEl: null | HTMLElement;
  isMenuOpen: boolean;
  handleMenuClose: () => void;
}
export default function BigScreenMenu({
  menuId,
  anchorEl,
  isMenuOpen,
  handleMenuClose,
}: BigScreenMenuProps) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isPending, mutate } = useMutation({
    mutationFn: logoutRequest,
    onSuccess: () => {
      tokenStore.clear();
      queryClient.setQueryData(["me"], null);
      queryClient.setQueryData(["me", "author"], null);
      handleMenuClose();
      localStorage.removeItem("hasSession");
    },
  });

  const handleLogout = () => {
    mutate();
  };

  const handleProfileNavigate = async () => {
    await navigate("/user/profile");
    handleMenuClose();
  };
  const handleProfileFollowedAuthors = async () => {
    await navigate("/user/followed-authors");
    handleMenuClose();
  };
  const handleProfileLikeBlogs = async () => {
    await navigate("/user/like-blogs");
    handleMenuClose();
  };
  return (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: 50,
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem
        disabled={isPending}
        onClick={handleProfileNavigate}
        sx={{
          display: "flex",
          placeContent: "flex-start",
        }}
      >
        Profile
      </MenuItem>
      <MenuItem
        disabled={isPending}
        onClick={handleProfileFollowedAuthors}
        sx={{
          display: "flex",
          placeContent: "flex-start",
        }}
      >
        Followed Authors
      </MenuItem>
      <MenuItem
        disabled={isPending}
        onClick={handleProfileLikeBlogs}
        sx={{
          display: "flex",
          placeContent: "flex-start",
        }}
      >
        Liked Blogs
      </MenuItem>
      <MenuItem
        disabled={isPending}
        onClick={handleLogout}
        sx={{
          display: "flex",
          placeContent: "flex-start",
        }}
      >
        {isPending ? <CircularProgress size="20px" /> : "Logout"}
      </MenuItem>
    </Menu>
  );
}
