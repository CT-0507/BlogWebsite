import { appName } from "@/config/const";
import Typography from "@mui/material/Typography";

export default function AboutPage() {
  const title = `${appName}  | About`;
  return (
    <>
      <title>{title}</title>
      <meta name="description" content="Information about this page." />
      <Typography>This is about page</Typography>
    </>
  );
}
