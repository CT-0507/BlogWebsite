import { appName } from "@/config/const";
import Typography from "@mui/material/Typography";

export default function UserHome() {
  const title = `${appName}  | My profile`;
  return (
    <>
      <title>{title}</title>
      <Typography component="h1">This is user home page</Typography>
    </>
  );
}
