import Grid from "@mui/material/Grid";
import Box from "@mui/material/Box";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Avatar from "@mui/material/Avatar";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import { useQuery } from "@tanstack/react-query";
import { getAuthorProfileRequest } from "@/api/authorApi";
import { relativeTime } from "@/utils/timeUtils";
import type { Author } from "@/types/types";
import FollowSection from "./FollowSection";
import { useAuth } from "@/hooks/useAuth";
import NotLoginFollowButton from "./NotLoginFollowButton";
import FollowersSection from "./FollowerSection";

interface AuthorProfileColumnProps {
  slug: string;
}

export default function AuthorProfileColumn({
  slug,
}: AuthorProfileColumnProps) {
  const isAuthenticated = useAuth();
  const { data, isLoading } = useQuery({
    queryKey: ["author", slug],
    queryFn: () => getAuthorProfileRequest(slug),
  });

  console.log(data);
  const author = data as Author;

  return (
    <Grid size={{ xs: 12, md: 3 }}>
      <Box sx={{ position: "sticky", top: 24 }}>
        <Card sx={{ borderRadius: 4 }}>
          {!isLoading && (
            <CardContent sx={{ textAlign: "center" }}>
              <Avatar
                //   src="https://i.pravatar.cc/150"
                sx={{ width: 80, height: 80, margin: "0 auto", mb: 2 }}
              />
              <Typography variant="h6">{author.displayName}</Typography>
              <Typography variant="body2" color="text.secondary">
                Frontend Engineer & Technical Writer
              </Typography>
              <Divider sx={{ my: 2 }} />
              {isAuthenticated ? (
                <FollowSection author={author} />
              ) : (
                <NotLoginFollowButton />
              )}
              <FollowersSection slug={slug} />
              <Divider sx={{ my: 2 }} />
              <Typography variant="body2">
                {author.bio}
                Passionate about React, UI architecture, and building scalable
                web apps.
              </Typography>
              <Divider sx={{ my: 2 }} />
              <Typography variant="body2">
                Follower count: {author.followerCount}
              </Typography>
              <Typography variant="body2">
                Blog count: {author.blogCount}
              </Typography>
              <Divider sx={{ my: 2 }} />
              <Typography variant="body2">
                Joined at: {relativeTime(author.createdAt)} ago
              </Typography>
            </CardContent>
          )}
        </Card>
      </Box>
    </Grid>
  );
}
