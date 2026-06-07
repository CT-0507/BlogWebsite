import useMediaQuery from "@mui/material/useMediaQuery";
import { useTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import MenuIcon from "@mui/icons-material/Menu";
import ChevronLeftIcon from "@mui/icons-material/ChevronLeft";
import DashboardIcon from "@mui/icons-material/Dashboard";
import SettingsIcon from "@mui/icons-material/Settings";
import ArticleIcon from "@mui/icons-material/Article";
import LogoutIcon from "@mui/icons-material/Logout";
import { useState } from "react";
import { Outlet, useNavigate } from "react-router-dom";
import NavBreadcrumbs from "@/components/Navigation/NavBreadcrumbs";
import BackButton from "@/components/Navigation/BackButton";
import Stack from "@mui/material/Stack";

const navItems = [
  { label: "Dashboard", icon: <DashboardIcon />, to: "dashboard" },
  { label: "My Blogs", icon: <ArticleIcon />, to: "my-blogs" },
  { label: "Settings", icon: <SettingsIcon />, to: "dashboard" },
  { label: "Logout", icon: <LogoutIcon />, to: "dashboard" },
];

export default function ResponsiveSidebarLayout() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));
  const navigate = useNavigate();

  const [open, setOpen] = useState(!isMobile);

  const toggleSidebar = () => {
    setOpen((prev) => !prev);
  };

  // Desktop widths
  const expandedWidth = "15%";
  const collapsedWidth = "72px";

  return (
    <Box
      id="dashboard-layout"
      sx={{ display: "flex", flex: 1, height: "100%" }}
    >
      {/* MOBILE FLOAT BUTTON */}
      {isMobile && (
        <IconButton
          onClick={toggleSidebar}
          sx={{
            position: "absolute",
            top: 16,
            left: 16,
            zIndex: 2000,
            bgcolor: "background.paper",
            boxShadow: 3,
          }}
        >
          {open ? <ChevronLeftIcon /> : <MenuIcon />}
        </IconButton>
      )}

      {/* SIDEBAR */}
      <Box
        sx={{
          width: isMobile
            ? open
              ? "100%"
              : 0
            : open
              ? expandedWidth
              : collapsedWidth,
          transition: "width 0.3s ease",
          position: "relative",
        }}
      >
        <Paper
          elevation={0}
          sx={{
            position: "sticky",
            top: 0, // sticks INSIDE parent container only
            height: "100%",
            overflow: "hidden",
            display: "flex",
            flexDirection: "column",
            width: "100%",
            borderRadius: 0,
          }}
        >
          {/* DESKTOP HEADER */}
          {!isMobile && (
            <>
              <Box
                sx={{
                  height: 64,
                  px: 1,
                  display: "flex",
                  alignItems: "center",
                  justifyContent: open ? "space-between" : "center",
                }}
              >
                {open && (
                  <Typography variant="h6" noWrap>
                    My App
                  </Typography>
                )}

                <IconButton onClick={toggleSidebar}>
                  {open ? <ChevronLeftIcon /> : <MenuIcon />}
                </IconButton>
              </Box>

              <Divider />
            </>
          )}

          {/* NAVIGATION */}
          <List sx={{ mt: isMobile ? 8 : 0 }}>
            {navItems.map((item) => (
              <ListItemButton
                key={item.label}
                sx={{
                  minHeight: 56,
                  justifyContent: open ? "initial" : "center",
                  px: 2.5,
                }}
                onClick={() => navigate(item.to)}
              >
                <ListItemIcon
                  sx={{
                    minWidth: 0,
                    mr: open ? 2 : "auto",
                    justifyContent: "center",
                  }}
                >
                  {item.icon}
                </ListItemIcon>

                {open && <ListItemText primary={item.label} />}
              </ListItemButton>
            ))}
          </List>
        </Paper>
      </Box>
      <Box id="content" sx={{ p: 1, flex: 1 }}>
        <Stack
          direction="row"
          sx={{ display: "flex", alignItems: "center" }}
          spacing={1}
        >
          <BackButton />
          <NavBreadcrumbs hiddenSegments={["view", "edit"]} />
        </Stack>
        <Outlet />
      </Box>
    </Box>
  );
}
