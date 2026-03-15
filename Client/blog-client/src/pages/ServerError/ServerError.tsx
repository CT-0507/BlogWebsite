import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import { Link } from "react-router-dom";

interface ServerErrorProps {
  status: number;
  statusText: string;
}

export default function ServerError({ status, statusText }: ServerErrorProps) {
  return (
    <Container sx={{ textAlign: "center", mt: 10 }}>
      <Typography variant="h2">500</Typography>

      <Typography variant="h5" sx={{ mt: 2 }}>
        Something went wrong
      </Typography>

      <Typography sx={{ mt: 2 }}>{status}</Typography>

      <Typography sx={{ mt: 2 }}>
        Eror Data
        {statusText}
      </Typography>

      <Button variant="contained" component={Link} to="/" sx={{ mt: 4 }}>
        Go Home
      </Button>
    </Container>
  );
}
