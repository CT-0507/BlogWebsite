import Box from "@mui/material/Box";
import BlogList from "./BlogList";
import { appName } from "@/config/const";

export default function Home() {
  const title = `${appName ?? ""} | Home`;
  return (
    <>
      <title>{title}</title>
      <meta
        name="description"
        content="Explore our latest blog posts, expert insights, practical guides, and industry updates. Discover valuable content to help you stay informed and inspired."
      />
      <Box>
        <BlogList />
      </Box>
    </>
  );
}
