import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import { Link } from "react-router-dom";
import { appName } from "@/config/const";

export default function Unauthorized() {
  const title = `${appName}  | Unauthorized`;
  return (
    <Container maxWidth="sm" sx={{ textAlign: "center", mt: 10 }}>
      <title>{title}</title>
      <meta
        name="description"
        content="You don't have permission to view this page"
      />
      <Stack spacing={3} alignItems="center">
        <Typography variant="h2" fontWeight="bold">
          401
        </Typography>

        <Typography variant="h5">Unauthorized</Typography>

        <Typography color="text.secondary">
          You must log in to access this page.
        </Typography>

        <Button variant="contained" component={Link} to="/account?action=login">
          Go to Login
        </Button>
      </Stack>
    </Container>
  );
}
