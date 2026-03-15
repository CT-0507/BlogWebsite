import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import { Link } from "react-router-dom";

export default function NotFound() {
  return (
    <Container maxWidth="sm" sx={{ textAlign: "center", mt: 10 }}>
      <Stack spacing={3} alignItems="center">
        <Typography variant="h2" fontWeight="bold">
          404
        </Typography>

        <Typography variant="h5">Page Not Found</Typography>

        <Typography color="text.secondary">
          The page you are looking for does not exist.
        </Typography>

        <Button variant="contained" component={Link} to="/">
          Go Home
        </Button>
      </Stack>
    </Container>
  );
}
