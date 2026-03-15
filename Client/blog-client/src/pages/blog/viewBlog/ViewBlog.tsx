import Container from "@mui/material/Container";
import Link from "@mui/material/Link";
import { Link as RouterLink, useNavigate } from "react-router-dom";
import Typography from "@mui/material/Typography";
import { useParams } from "react-router-dom";
import BlogSection from "./BlogSection";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Divider from "@mui/material/Divider";
import SocialShareDial from "@/components/SocialShareSpeedDial/SocialShareSpeedDial";

export default function ViewPost() {
  const navigate = useNavigate();
  const { id } = useParams();

  if (!id)
    return (
      <>
        <Typography>Id not found</Typography>
        <Link component={RouterLink} to="/">
          Go back
        </Link>
      </>
    );

  const handleGoBack = () => {
    navigate(-1);
  };

  return (
    <Container maxWidth="md" sx={{ mt: 5 }}>
      {/* Thumbnail */}
      <BlogSection id={id} />
      <Divider />
      <Box
        id="footer-action"
        mt={1}
        sx={{
          display: "flex",
          justifyContent: "flex-end",
        }}
      >
        <Button onClick={handleGoBack} variant="contained">
          Go back
        </Button>
      </Box>

      <SocialShareDial />
    </Container>
  );
}
