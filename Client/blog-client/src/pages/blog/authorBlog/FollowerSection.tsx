import { getAuthorFollowersRequest } from "@/api/authorApi";
import { truncate } from "@/utils/textUtils";
import Avatar from "@mui/material/Avatar";
import AvatarGroup from "@mui/material/AvatarGroup";
import Box from "@mui/material/Box";
import Skeleton from "@mui/material/Skeleton";
import Tooltip from "@mui/material/Tooltip";
import Typography from "@mui/material/Typography";
import { useQuery } from "@tanstack/react-query";

// type Follower = {
//   id: number;
//   avatarUrl: string;
// };

interface FollowersSectionProps {
  slug: string;
}

export default function FollowersSection({ slug }: FollowersSectionProps) {
  const { data, isLoading } = useQuery({
    queryKey: ["followers", slug],
    queryFn: () => getAuthorFollowersRequest(slug),
    staleTime: 1000 * 60 * 5,
  });

  if (isLoading) {
    return (
      <Box display="flex" alignItems="center" gap={1}>
        <Skeleton variant="circular" width={32} height={32} />
        <Skeleton variant="circular" width={32} height={32} />
        <Skeleton variant="text" width={80} />
      </Box>
    );
  }

  const followers = data ?? [];
  const display = followers.slice(0, 5);

  return (
    <Box display="flex" alignItems="center" gap={1}>
      <AvatarGroup max={5}>
        {display.map((follower: string) => (
          <Tooltip title={truncate(follower, 20)} key={follower}>
            <Avatar
              // src={follower.avatarUrl}
              alt={`Follower ${follower}`}
            />
          </Tooltip>
        ))}
      </AvatarGroup>

      <Typography variant="body2" color="text.secondary">
        {followers.length} followers
      </Typography>
    </Box>
  );
}
