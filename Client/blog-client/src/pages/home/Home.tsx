import Box from "@mui/material/Box";
import BlogList from "./BlogList";
import Typography from "@mui/material/Typography";

export default function Home() {
  return (
    <>
      <Box>
        <Box
          id="home-header"
          sx={{
            display: "flex",
            justifyContent: "center",
            mt: 1,
          }}
        >
          <Typography variant="h3" display="inline">
            This is home page
          </Typography>
        </Box>
        <BlogList />
      </Box>
    </>
  );
}
