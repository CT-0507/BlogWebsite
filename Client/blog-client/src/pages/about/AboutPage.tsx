import { appName } from "@/config/const";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import Markdown from "react-markdown";
import about from "@/assets/README.md?raw";

export default function AboutPage() {
  const title = `${appName}  | About`;
  return (
    <Box
      sx={{
        flex: 1,
        p: 1,
        px: 2,
      }}
    >
      <title>{title}</title>
      <meta name="description" content="Information about this page." />
      <Typography variant="h3" align="center">
        About this Project
      </Typography>
      <Divider />
      <Box sx={{ px: 2 }}>
        <Markdown>{about}</Markdown>
      </Box>
    </Box>
  );
}
