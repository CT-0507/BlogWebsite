import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import CreateAuthorForm from "./CreateAuthorForm";

export default function CreateAuthorPage() {
  return (
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
      <Typography variant="h3">Become an Author today</Typography>
      <Divider />
      <CreateAuthorForm />
    </Box>
  );
}
