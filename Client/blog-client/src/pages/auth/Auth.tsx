import CardContent from "@mui/material/CardContent";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardMedia from "@mui/material/CardMedia";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import { useState } from "react";
import LoginForm from "./LoginForm";
import SignupForm from "./SignupForm";

interface TabPanelProps {
  children?: React.ReactNode;
  dir?: string;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`full-width-tabpanel-${index}`}
      aria-labelledby={`full-width-tab-${index}`}
      {...other}
    >
      {value === index && <Box sx={{ p: 3 }}>{children}</Box>}
    </div>
  );
}

function a11yProps(tab: number) {
  return {
    id: `full-width-tab-${tab}`,
    "aria-controls": `full-width-tabpanel-${tab}`,
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
  const [currentTab, setCurrentTab] = useState<number>(0);
  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    setCurrentTab(newValue);
  };
  return (
    <Box
      sx={{
        minWidth: "500px",
        height: "100%",
        width: "50%",
        m: "auto",
      }}
    >
      <Card
        sx={{
          width: "100%",
          minHeight: "500px",
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
            <Button onClick={(e) => handleTabChange(e, 1)}>
              Don't have an account yet
            </Button>
          ) : (
            <Button onClick={(e) => handleTabChange(e, 0)}>
              Already have an account?
            </Button>
          )}
        </Box>
      </Card>
    </Box>
  );
}
