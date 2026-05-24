import { CardContent, CardHeader, CircularProgress } from "@mui/material";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Card from "@mui/material/Card";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import { useGetAuthorDashboardData } from "@/hooks/useMetrics";

export default function AuthorDashboard() {
  const { data, isLoading } = useGetAuthorDashboardData();
  return (
    <Box>
      <Typography>This is author dashboard</Typography>
      <Stack>
        <Paper>
          <Card>
            <CardHeader>
              <Typography variant="h3">Total Views</Typography>
            </CardHeader>
            <CardContent>
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
        </Paper>
        <Paper>
          <Card>
            <CardHeader>
              <Typography variant="h3">Total Likes</Typography>
            </CardHeader>
            <CardContent>
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
        </Paper>
        <Paper>
          <Card>
            <CardHeader>
              <Typography variant="h3">Total Dislikes</Typography>
            </CardHeader>
            <CardContent>
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
        </Paper>
      </Stack>
    </Box>
  );
}
