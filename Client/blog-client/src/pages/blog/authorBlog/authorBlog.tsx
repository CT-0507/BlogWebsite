import {
  Box,
  Card,
  CardContent,
  Typography,
  List,
  ListItem,
  ListItemText,
  Divider,
} from "@mui/material";
import Grid from "@mui/material/Grid";
import AuthorProfileColumn from "./AuthorProfileColumn";
import { useParams } from "react-router-dom";
import AuthorBlogsColumn from "./AuthorBlogsColumn";

export default function AuthorBlogPage() {
  const { slug } = useParams();
  if (!slug) {
    return <h1>Author slug not found</h1>;
  }
  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Grid container spacing={3}>
        {/* Left Column - Author Info */}
        <AuthorProfileColumn slug={slug} />

        {/* Middle Column - Blogs */}
        <AuthorBlogsColumn slug={slug} />

        {/* Right Column - Relevant Info */}
        <Grid size={{ xs: 12, md: 3 }}>
          <Box sx={{ position: "sticky", top: 24 }}>
            <Card sx={{ borderRadius: 4 }}>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Related Topics
                </Typography>
                <List>
                  <ListItem>
                    <ListItemText primary="React Performance" />
                  </ListItem>
                  <ListItem>
                    <ListItemText primary="State Management" />
                  </ListItem>
                  <ListItem>
                    <ListItemText primary="UI/UX Patterns" />
                  </ListItem>
                  <ListItem>
                    <ListItemText primary="TypeScript with React" />
                  </ListItem>
                </List>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2" color="text.secondary">
                  Subscribe to get updates on new articles and tutorials.
                </Typography>
              </CardContent>
            </Card>
          </Box>
        </Grid>
      </Grid>
    </Box>
  );
}
