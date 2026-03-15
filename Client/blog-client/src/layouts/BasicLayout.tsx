import { tokenStore } from "@/api/store/tokenStore";
import { useAuth } from "@/hooks/useAuth";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import IconButton from "@mui/material/IconButton";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import AccountCircle from "@mui/icons-material/AccountCircle";
import MoreIcon from "@mui/icons-material/MoreVert";
import { Link, Outlet, useNavigate } from "react-router-dom";
import { useState } from "react";
import BorderColor from "@mui/icons-material/BorderColor";
import Button from "@mui/material/Button";
import Snackbar from "@mui/material/Snackbar";
import CloseIcon from "@mui/icons-material/Close";
import { useAuthSSE } from "@/hooks/useSSECacheBridge";
import NotificationMenu from "./NotificationMenu";
import MobileMenu from "./MobileMenu";
import BigScreenMenu from "./BigScreenMenu";

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

  useAuthSSE(
    isAuthenticated ? tokenStore.get() : null,
    [],
    ["blog_created_admin"],
    setSnackbarOpen
  );

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
