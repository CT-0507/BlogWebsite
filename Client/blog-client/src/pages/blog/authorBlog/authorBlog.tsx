import {
  Box,
  Card,
  CardContent,
  Typography,
  Avatar,
  List,
  ListItem,
  ListItemText,
  Divider,
} from "@mui/material";
import Grid from "@mui/material/Grid";

const blogs = [
  {
    id: 1,
    title: "Understanding React Hooks",
    excerpt:
      "Learn how React Hooks simplify state management and side effects.",
  },
  {
    id: 2,
    title: "Material UI Layout Tips",
    excerpt:
      "Best practices for building responsive layouts using Material UI.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
  {
    id: 3,
    title: "Scaling Frontend Architecture",
    excerpt: "Strategies to keep large React applications maintainable.",
  },
];

export default function AuthorBlogPage() {
  return (
    <Box sx={{ flexGrow: 1, p: 3 }}>
      <Grid container spacing={3}>
        {/* Left Column - Author Info */}
        <Grid size={{ xs: 12, md: 3 }}>
          <Box sx={{ position: "sticky", top: 24 }}>
            <Card sx={{ borderRadius: 4 }}>
              <CardContent sx={{ textAlign: "center" }}>
                <Avatar
                  src="https://i.pravatar.cc/150"
                  sx={{ width: 80, height: 80, margin: "0 auto", mb: 2 }}
                />
                <Typography variant="h6">John Doe</Typography>
                <Typography variant="body2" color="text.secondary">
                  Frontend Engineer & Technical Writer
                </Typography>
                <Divider sx={{ my: 2 }} />
                <Typography variant="body2">
                  Passionate about React, UI architecture, and building scalable
                  web apps.
                </Typography>
              </CardContent>
            </Card>
          </Box>
        </Grid>

        {/* Middle Column - Blogs */}
        <Grid size={{ xs: 12, md: 6 }}>
          <Grid container spacing={2}>
            {blogs.map((blog) => (
              <Grid size={{ xs: 12 }} key={blog.id}>
                <Card sx={{ borderRadius: 4 }}>
                  <CardContent>
                    <Typography variant="h6" gutterBottom>
                      {blog.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {blog.excerpt}
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        </Grid>

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
