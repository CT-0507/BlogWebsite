import { logoutRequest } from "@/api/auth";
import { tokenStore } from "@/api/store/tokenStore";
import { useAuth } from "@/hooks/useAuth";
import AppBar from "@mui/material/AppBar";
import Badge from "@mui/material/Badge";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import AccountCircle from "@mui/icons-material/AccountCircle";
import NotificationsIcon from "@mui/icons-material/Notifications";
import MoreIcon from "@mui/icons-material/MoreVert";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, Outlet, useNavigate } from "react-router-dom";
import { useEffect, useMemo, useState } from "react";
import CircularProgress from "@mui/material/CircularProgress";
import BorderColor from "@mui/icons-material/BorderColor";
import Button from "@mui/material/Button";
import { getNotifications } from "@/api/userApi";
import Snackbar from "@mui/material/Snackbar";
import CloseIcon from "@mui/icons-material/Close";
import { useAuthSSE } from "@/hooks/useSSECacheBridge";

interface BigScreenMenuProps {
  menuId: string;
  anchorEl: null | HTMLElement;
  isMenuOpen: boolean;
  handleMenuClose: () => void;
}
function BigScreenMenu({
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
        vertical: "top",
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
          width: "80px",
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
          width: "80px",
          display: "flex",
          placeContent: "center",
        }}
      >
        {isPending ? <CircularProgress size="20px" /> : "Logout"}
      </MenuItem>
    </Menu>
  );
}

interface MobileMenuProps {
  mobileMenuId: string;
  mobileMoreAnchorEl: null | HTMLElement;
  isMobileMenuOpen: boolean;
  handleMobileMenuClose: (event: React.MouseEvent<HTMLElement>) => void;
  handleProfileMenuOpen: (event: React.MouseEvent<HTMLElement>) => void;
}

function MobileMenu({
  mobileMenuId,
  mobileMoreAnchorEl,
  isMobileMenuOpen,
  handleMobileMenuClose,
  handleProfileMenuOpen,
}: MobileMenuProps) {
  const { user } = useAuth();
  return (
    <Menu
      sx={{ zIndex: 99 }}
      anchorEl={mobileMoreAnchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={mobileMenuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMobileMenuOpen}
      onClose={handleMobileMenuClose}
    >
      <MenuItem>
        <IconButton
          size="large"
          aria-label="show 17 new notifications"
          color="inherit"
        >
          <Badge badgeContent={17} color="error">
            <NotificationsIcon />
          </Badge>
        </IconButton>
        <p>Notifications</p>
      </MenuItem>
      {user && (
        <MenuItem onClick={handleProfileMenuOpen}>
          <IconButton
            size="large"
            aria-label="account of current user"
            aria-controls="primary-search-account-menu"
            aria-haspopup="true"
            color="inherit"
          >
            <AccountCircle />
          </IconButton>
          <p>Profile</p>
        </MenuItem>
      )}
    </Menu>
  );
}

function NotificationMenu() {
  const { data, isLoading } = useQuery({
    queryKey: ["notifications"],
    queryFn: getNotifications,
    retry: false,
    refetchInterval: 30 * 60 * 1000,
  });

  const nofiticationNumber = useMemo(() => {
    if (!isLoading && data && Array.isArray(data)) {
      return data.length;
    }
    return 0;
  }, [data, isLoading]);

  const handleShowNotifications = () => {};

  useEffect(() => {
    document.title =
      document.title +
      (nofiticationNumber !== 0 ? `(${nofiticationNumber})` : "");
  }, [nofiticationNumber]);

  return (
    <IconButton
      size="large"
      aria-label="show 17 new notifications"
      color="inherit"
      onClick={handleShowNotifications}
    >
      <Badge badgeContent={nofiticationNumber} color="error">
        <NotificationsIcon />
      </Badge>
    </IconButton>
  );
}

export default function BasicLayout() {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [mobileMoreAnchorEl, setMobileMoreAnchorEl] =
    useState<null | HTMLElement>(null);
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const isMenuOpen = Boolean(anchorEl);
  const isMobileMenuOpen = Boolean(mobileMoreAnchorEl);

  const menuId = "primary-search-account-menu";
  const mobileMenuId = "primary-search-account-menu-mobile";

  console.log("Rerender");
  console.log(isAuthenticated);
  useAuthSSE(isAuthenticated ? tokenStore.get() : null, undefined, [
    "blog_created_admin",
  ]);

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMobileMenuClose = () => {
    setMobileMoreAnchorEl(null);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    handleMobileMenuClose();
  };

  const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMobileMoreAnchorEl(event.currentTarget);
  };

  const handleLogoClick = () => {
    navigate("/dashboard");
  };

  const handleCloseSnackBar = () => {
    setSnackbarOpen(false);
  };

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        minHeight: "100vh",
      }}
      component="main"
    >
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={handleCloseSnackBar}
        message="Note archived"
        action={
          <>
            <Typography>You have notification</Typography>
            <IconButton
              size="small"
              aria-label="close"
              color="inherit"
              onClick={handleCloseSnackBar}
            >
              <CloseIcon fontSize="small" />
            </IconButton>
          </>
        }
      />
      {/* Header */}
      <Box display="contents" sx={{ flexGrow: 1 }}>
        <AppBar position="static">
          <Toolbar>
            {/* Logo placeholder */}
            <Box
              component="img"
              src="/logo.png"
              alt="Logo"
              sx={{ height: 40, mr: 2, cursor: "pointer" }}
              onClick={handleLogoClick}
            />
            <Typography variant="h6" component="div">
              My App Name
            </Typography>
            <Box sx={{ flexGrow: 1 }} />
            <Box sx={{ display: { xs: "none", md: "flex" } }}>
              {isAuthenticated ? (
                <>
                  <NotificationMenu />
                  <Button
                    component={Link}
                    to="/blog/publish"
                    size="large"
                    aria-label="account of current user"
                    aria-haspopup="true"
                    color="info"
                    title="Publish new blog"
                    variant="contained"
                    sx={{
                      mx: 1,
                    }}
                  >
                    <BorderColor />
                    <Typography ml={2}>Publish new blog</Typography>
                  </Button>
                  <IconButton
                    size="large"
                    edge="end"
                    aria-label="account of current user"
                    aria-haspopup="true"
                    onClick={handleProfileMenuOpen}
                    color="inherit"
                    title="Profile"
                  >
                    <AccountCircle />
                  </IconButton>
                </>
              ) : (
                <>
                  <Button component={Link} to="/account" color="info">
                    Account
                  </Button>
                </>
              )}
            </Box>
            <Box sx={{ display: { xs: "flex", md: "none" } }}>
              <IconButton
                size="large"
                aria-label="show more"
                aria-controls={mobileMenuId}
                aria-haspopup="true"
                onClick={handleMobileMenuOpen}
                color="inherit"
              >
                <MoreIcon />
              </IconButton>
            </Box>
          </Toolbar>
        </AppBar>
        <BigScreenMenu
          menuId={menuId}
          anchorEl={anchorEl}
          isMenuOpen={isMenuOpen}
          handleMenuClose={handleMenuClose}
        />
        <MobileMenu
          mobileMenuId={mobileMenuId}
          mobileMoreAnchorEl={mobileMoreAnchorEl}
          isMobileMenuOpen={isMobileMenuOpen}
          handleMobileMenuClose={handleMobileMenuClose}
          handleProfileMenuOpen={handleProfileMenuOpen}
        />
      </Box>

      {/* Main Content */}
      <Box
        component="main"
        sx={{
          flex: 1, // pushes footer down
          py: 1,
          justifyContent: "center",
        }}
      >
        <Outlet />
      </Box>

      {/* Footer */}
      <Box
        component="footer"
        sx={{
          py: 2,
          px: 2,
          mt: "auto",
          backgroundColor: (theme) => theme.palette.grey[200],
          textAlign: "center",
        }}
      >
        <Typography variant="body2" color="text.secondary">
          © {new Date().getFullYear()} My App. All rights reserved.
        </Typography>
      </Box>
    </Box>
  );
}
