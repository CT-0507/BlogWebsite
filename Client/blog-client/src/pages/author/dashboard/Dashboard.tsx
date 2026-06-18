import Box from "@mui/material/Box";
import CardContent from "@mui/material/CardContent";
import CircularProgress from "@mui/material/CircularProgress";
import Typography from "@mui/material/Typography";
import Card from "@mui/material/Card";
import Stack from "@mui/material/Stack";
import { useGetAuthorDashboardData } from "@/hooks/useMetrics";
import { appName } from "@/config/const";

export default function AuthorDashboard() {
  const { data, isLoading } = useGetAuthorDashboardData();
  const title = `${appName} | Author dashboard`;
  return (
    <Box>
      <title>{title}</title>
      <Typography>This is author dashboard</Typography>
      <Stack direction="row">
        <Card sx={{ mr: 4, p: 2 }}>
          <CardContent>
            <Typography variant="h5" fontWeight="600" sx={{ mb: 2 }}>
              Total Views
            </Typography>
            <Stack>
              <Typography>
                Today views{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.viewsMetrics.todayViews
                )}
              </Typography>
              <Typography>
                Yesterday views{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.viewsMetrics.yesterdayViews
                )}
              </Typography>
              <Typography>
                This week views{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.viewsMetrics.thisWeekViews
                )}
              </Typography>
              <Typography>
                Last week views{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.viewsMetrics.lastWeekViews
                )}
              </Typography>
            </Stack>
          </CardContent>
        </Card>
        <Card sx={{ mr: 4, p: 2 }}>
          <CardContent>
            <Typography variant="h5" fontWeight="600" sx={{ mb: 2 }}>
              Total Likes
            </Typography>
            <Stack>
              <Typography>
                Today likes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.todayLikes
                )}
              </Typography>
              <Typography>
                Yesterday likes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.yesterdayLikes
                )}
              </Typography>
              <Typography>
                This week likes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.thisWeekLikes
                )}
              </Typography>
              <Typography>
                Last week likes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.lastWeekLikes
                )}
              </Typography>
            </Stack>
          </CardContent>
        </Card>
        <Card sx={{ mr: 4, p: 2 }}>
          <CardContent>
            <Typography variant="h5" fontWeight="600" sx={{ mb: 2 }}>
              Total Dislikes
            </Typography>
            <Stack>
              <Typography>
                Today dislikes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.todayDislikes
                )}
              </Typography>
              <Typography>
                Yesterday dislikes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.yesterdayDislikes
                )}
              </Typography>
              <Typography>
                This week dislikes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.thisWeekDislikes
                )}
              </Typography>
              <Typography>
                Last week dislikes{" "}
                {isLoading ? (
                  <CircularProgress />
                ) : (
                  data?.reactionMetrics.lastWeekDislikes
                )}
              </Typography>
            </Stack>
          </CardContent>
        </Card>
      </Stack>
    </Box>
  );
}
