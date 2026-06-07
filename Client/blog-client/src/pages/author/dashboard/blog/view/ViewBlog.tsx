import ViewBlog from "@/pages/blog/viewBlog/ViewBlog";
import Box from "@mui/material/Box";
import ChartSection from "../components/ChartSection";
import { useState } from "react";
import { useParams } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";
import { useBlogBySlug } from "@/hooks/useBlogBySlug";
import { useGetAuthorDashboardBlogMetrics } from "@/hooks/useMetrics";
import CircularProgress from "@mui/material/CircularProgress";
import { formatDayShort, formatWeekLabel } from "@/utils/timeUtils";

// const dailyViews = [
//   { date: "2026/05/29", views: 142 },
//   { date: "2026/05/30", views: 158 },
//   { date: "2026/05/31", views: 151 },
//   { date: "2026/06/01", views: 176 },
// ];

// const weeklyViews = [
//   { weekStart: "2026/05/04", views: 920 },
//   { weekStart: "2026/05/11", views: 1050 },
//   { weekStart: "2026/05/18", views: 1125 },
//   { weekStart: "2026/05/25", views: 1280 },
// ];
export default function ViewDashboardBlogPage() {
  const [weeks, setWeeks] = useState(4);
  const [dates, setDates] = useState(4);
  const { slug } = useParams();
  const { isAuthenticated } = useAuth();

  const { data: blog, isLoading } = useBlogBySlug(isAuthenticated, slug);

  const { data: dailyViewData } = useGetAuthorDashboardBlogMetrics({
    blogID: blog ? blog.blogID : 0,
    viewType: "days",
    resultLength: dates,
    enabled: !isLoading,
  });
  const { data: weeklyViewData } = useGetAuthorDashboardBlogMetrics({
    blogID: blog ? blog.blogID : 0,
    viewType: "weeks",
    resultLength: weeks,
    enabled: !isLoading,
  });

  if (isLoading) {
    return <CircularProgress />;
  }

  return (
    <Box>
      <ChartSection
        title="Daily Views"
        data={
          dailyViewData
            ? dailyViewData.data.map(
                (item: { date: string; views: number }) => ({
                  ...item,
                  date: formatDayShort(item.date),
                }),
              )
            : []
        }
        periodField="date"
        valueField="views"
        comparisonLabel="vs yesterday"
        sliderLabel="D"
        range={dates}
        onRangeChange={setDates}
        maxRange={7}
      />
      <ChartSection
        title="Weekly Views"
        data={
          weeklyViewData
            ? weeklyViewData.data.map(
                (item: { weekStart: string; views: number }) => ({
                  ...item,
                  weekStart: formatWeekLabel(item.weekStart),
                }),
              )
            : []
        }
        periodField="weekStart"
        valueField="views"
        comparisonLabel="vs last week"
        sliderLabel="D"
        range={weeks}
        onRangeChange={setWeeks}
        maxRange={7}
      />
      <ViewBlog blog={blog} isLoading={isLoading} />
    </Box>
  );
}
