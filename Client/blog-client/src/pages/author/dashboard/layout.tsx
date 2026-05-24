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
import PeopleIcon from "@mui/icons-material/People";
import LogoutIcon from "@mui/icons-material/Logout";
import { useState } from "react";
import { Outlet } from "react-router-dom";

const navItems = [
  { label: "Dashboard", icon: <DashboardIcon /> },
  { label: "Users", icon: <PeopleIcon /> },
  { label: "Settings", icon: <SettingsIcon /> },
  { label: "Logout", icon: <LogoutIcon /> },
];

export default function ResponsiveSidebarLayout() {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const [open, setOpen] = useState(!isMobile);

  const toggleSidebar = () => {
    setOpen((prev) => !prev);
  };

  // Desktop widths
  const expandedWidth = "15%";
  const collapsedWidth = "72px";

  return (
    <Box sx={{ display: "flex" }}>
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
          flexShrink: 0,
          position: "relative",
        }}
      >
        <Paper
          elevation={2}
          sx={{
            position: "sticky",
            top: 0, // sticks INSIDE parent container only
            height: "100%",
            mt: -1,
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
      <Box sx={{ p: 1 }}>
        <Outlet />
      </Box>
    </Box>
  );
}
