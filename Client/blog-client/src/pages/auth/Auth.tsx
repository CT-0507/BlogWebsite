import CardContent from "@mui/material/CardContent";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardMedia from "@mui/material/CardMedia";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import { useMemo, useState } from "react";
import LoginForm from "./LoginForm";
import SignupForm from "./SignupForm";
import { Navigate, useLocation, useSearchParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import Grid from "@mui/material/Grid";
import { Divider } from "@mui/material";
import { appName } from "@/config/const";

interface TabPanelProps {
  children?: React.ReactNode;
  dir?: string;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;
  return (
    <Box
      sx={{
        display: "flex",
        width: value !== index ? "" : "100%",
        minHeight: "450px",
      }}
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index == 0 ? "login" : "signup"}`}
      aria-labelledby={`full-width-tab-${index == 0 ? "login" : "signup"}`}
      {...other}
    >
      {value === index && (
        <Box
          sx={{ display: "flex", flexDirection: "column", p: 1, width: "100%" }}
        >
          {children}
        </Box>
      )}
    </Box>
  );
}

function a11yProps(tab: number) {
  return {
    id: `${tab == 0 ? "login-panel" : "signup-panel"}`,
    "aria-controls": `tabpanel-${tab == 0 ? "login" : "signup"}`,
  };
}

function FormHeader() {
  return (
    <>
      <Typography component="h1">Welcome to our blog website</Typography>
    </>
  );
}

const tabs = ["Login", "Register"];
export default function Login() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { isAuthenticated } = useAuth();
  const location = useLocation();
  const from = location.state?.from?.pathname || "/";

  const action = useMemo(() => {
    const query = searchParams.get("action");
    switch (query) {
      case "register":
        return 1;
      default:
        return 0;
    }
  }, [searchParams]);

  const [currentTab, setCurrentTab] = useState<number>(action);
  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
    let action: string;
    switch (newValue) {
      case 0:
        action = "login";
        break;
      default:
        action = "register";
        break;
    }
    setSearchParams({ action });
  };
  if (isAuthenticated) {
    return <Navigate to={from} replace />;
  }
  const title = `${appName} | ${action === 0 ? "Login" : "Register"}`;
  return (
    <Grid
      container
      direction="row"
      justifyContent="space-between"
      spacing={2}
      sx={{}}
    >
      <title>{title}</title>
      <meta
        name="description"
        content={
          action === 0
            ? "Explore our latest blog posts, expert insights, practical guides, and industry updates. Discover valuable content to help you stay informed and inspired."
            : "Create an account to access exclusive content, personalize your experience, and stay updated with the latest posts and features."
        }
      />
      <Grid
        size={{ md: 3 }}
        sx={{
          display: { xs: "none", md: "flex" },
          justifyContent: "center",
          alignItems: "center",
          flexDirection: "column",
        }}
      >
        <Box display={action === 0 ? "block" : "none"}>
          Test user:
          <br />
          Username: user1
          <br />
          Password: Abc!2345
        </Box>
        <Divider />
      </Grid>

      <Grid
        size={{ xs: 12, md: 6 }}
        order={{ xs: 3, md: 2 }}
        display="flex"
        justifyContent="center"
      >
        <Card
          sx={{
            width: "100%",
            minHeight: "500px",
            maxWidth: "600px",
            display: "flex",
            flexDirection: "column",
          }}
        >
          <CardMedia
            sx={{
              height: "30px",
              p: 1,
            }}
          >
            <FormHeader />
          </CardMedia>
          <Box
            sx={{
              width: "100%",
            }}
          >
            <Tabs
              value={currentTab}
              onChange={handleTabChange}
              indicatorColor="secondary"
              textColor="inherit"
              variant="fullWidth"
              aria-label="tabs"
            >
              {tabs.map((label, index) => (
                <Tab
                  key={index}
                  sx={{ fontWeight: "bold", fontSize: 20 }}
                  label={label}
                  {...a11yProps(index)}
                />
              ))}
            </Tabs>
          </Box>
          <CardContent
            sx={{
              width: "100%",
              flex: 1,
              display: "flex",
            }}
          >
            <TabPanel value={currentTab} index={0}>
              <LoginForm />
            </TabPanel>
            <TabPanel value={currentTab} index={1}>
              <SignupForm />
            </TabPanel>
          </CardContent>
          <Box
            sx={{
              margin: "auto",
            }}
          >
            {currentTab == 0 ? (
              <Button
                aria-labelledby="go-to-login"
                onClick={(e) => handleTabChange(e, 1)}
              >
                Don't have an account yet
              </Button>
            ) : (
              <Button
                aria-labelledby="go-to-signup"
                onClick={(e) => handleTabChange(e, 0)}
              >
                Already have an account?
              </Button>
            )}
          </Box>
        </Card>
      </Grid>
      <Grid size={{ xs: 0, md: 3 }} order={{ xs: 1, md: 3 }}></Grid>
    </Grid>
  );
}
