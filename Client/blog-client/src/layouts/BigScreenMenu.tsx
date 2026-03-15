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
      handleMenuClose();
      navigate("/account");
    },
    onError: (error) => {
      console.log(error);
    },
  });

  const handleLogout = () => {
    mutate();
  };

  const handleProfileNavigate = async () => {
    await navigate("/user/profile");
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
          width: "100px",
          display: "flex",
          placeContent: "center",
        }}
      >
        Profile
      </MenuItem>
      <MenuItem
        disabled={isPending}
        onClick={handleLogout}
        sx={{
          width: "100px",
          display: "flex",
          placeContent: "center",
        }}
      >
        {isPending ? <CircularProgress size="20px" /> : "Logout"}
      </MenuItem>
    </Menu>
  );
}
