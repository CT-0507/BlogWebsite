import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import type { Blog } from "@/pages/home/BlogList";
import { getAuthorBlogsRequest } from "@/api/authorApi";
import { useQuery } from "@tanstack/react-query";
import { relativeTime } from "@/utils/timeUtils";
import { truncate } from "@/utils/textUtils";

interface AuthorBlogsColumnProps {
  slug: string;
}
export default function AuthorBlogsColumn({ slug }: AuthorBlogsColumnProps) {
  const { data, isLoading } = useQuery({
    queryKey: ["author_blogs", slug],
    queryFn: () => getAuthorBlogsRequest(slug),
    staleTime: Infinity,
  });
  const blogs = data as Blog[];
  return (
    <Grid size={{ xs: 12, md: 6 }}>
      {!isLoading && (
        <Grid container spacing={2}>
          {blogs?.map((blog) => (
            <Grid size={{ xs: 12 }} key={blog.blogID}>
              <Card sx={{ borderRadius: 4 }}>
                <CardContent>
                  <Typography variant="h6" gutterBottom>
                    {blog.title} |{" "}
                    <Typography sx={{ fontSize: "0.75rem" }} component="span">
                      {relativeTime(blog.createdAt)}
                    </Typography>
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {truncate(blog.content, 100)}
                  </Typography>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      )}
    </Grid>
  );
}
