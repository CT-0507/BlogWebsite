import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import EmailSection from "./EmailSection";
import BasicInfoSection from "./BasicInfoSection";
import PasswordSection from "./PasswordSection";

export default function Profile() {
  return (
    <>
      <Box
        sx={{
          width: "100%",
          maxWidth: {
            md: "80%",
          },
          mx: "auto", // center horizontally
          p: 1,
        }}
      >
        <Typography component="h1">This is profile page</Typography>
        <Divider />
        <EmailSection />
        <Divider />
        <BasicInfoSection />
        <Divider />
        <PasswordSection />
        <Divider />
      </Box>
    </>
  );
}
