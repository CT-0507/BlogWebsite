import useRankingBlogs from "@/hooks/useRankingBlogs";
import { getTypeValidValue } from "@/utils/mapper";
import Box from "@mui/material/Box";
import Pagination from "@mui/material/Pagination";
import Stack from "@mui/material/Stack";
import Tab from "@mui/material/Tab";
import Tabs from "@mui/material/Tabs";
import Typography from "@mui/material/Typography";
import { useState } from "react";

interface BlogListProps {
  type: string;
}
function TrendingBlogList({ type }: BlogListProps) {
  const [page, setPage] = useState(1);
  const limit = 5;

  const { data, isLoading } = useRankingBlogs(
    {
      sortBy: "rank",
      sortDir: "asc",
      limit: limit,
    },
    page,
    getTypeValidValue(type, ["allTime", "trending"], "trending"),
    type === "trending"
  );

  return (
    <div hidden={type !== "trending"}>
      <Typography variant="subtitle1" fontWeight={600} gutterBottom>
        Top trending blogs
      </Typography>

      {/* Desktop: pagination list */}
      <Box sx={{ display: { xs: "none", md: "block" } }}>
        <Stack spacing={1}>
          {!isLoading &&
            data?.blogs.map((item, i) => (
              <Box key={i} sx={{ height: 30, background: "#f0f0f0" }}>
                {item.title ?? "Deleted"}
              </Box>
            ))}
        </Stack>
      </Box>

      {/* Mobile: swipe carousel (simple horizontal scroll) */}
      <Box
        sx={{
          display: { xs: "flex", md: "none" },
          overflowX: "auto",
          scrollSnapType: "x mandatory",
          // display: "flex",
          gap: 1,
        }}
      >
        {!isLoading &&
          data?.blogs.map((item, i) => (
            <Box
              key={i}
              sx={{
                minWidth: "80%",
                height: 80,
                background: "#f0f0f0",
                scrollSnapAlign: "center",
                flexShrink: 0,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <Typography>{item.title ?? "Deleted"}</Typography>
              <Typography>{item.content ?? "Deleted"}</Typography>
            </Box>
          ))}
      </Box>
      <Pagination
        count={Math.ceil(data?.total ?? 0 / limit)}
        page={page}
        onChange={(_, value) => setPage(value)}
        size="small"
        sx={{ mt: 2 }}
      />
    </div>
  );
}

function AllTimeBlogList({ type }: BlogListProps) {
  const [page, setPage] = useState(1);
  const limit = 5;

  const { data, isLoading } = useRankingBlogs(
    {
      sortBy: "rank",
      sortDir: "asc",
      limit: limit,
    },
    page,
    getTypeValidValue(type, ["allTime", "trending"], "allTime"),
    type === "allTime"
  );

  return (
    <div hidden={type !== "allTime"}>
      <Typography variant="subtitle1" fontWeight={600} gutterBottom>
        Top all time blogs
      </Typography>

      {/* Desktop: pagination list */}
      <Box sx={{ display: { xs: "none", md: "block" } }}>
        <Stack spacing={1}>
          {!isLoading &&
            data?.blogs.map((item, i) => (
              <Box key={i} sx={{ height: 30, background: "#f0f0f0" }}>
                {item.title ?? "Deleted"}
              </Box>
            ))}
        </Stack>
      </Box>

      {/* Mobile: swipe carousel (simple horizontal scroll) */}
      <Box
        sx={{
          display: { xs: "flex", md: "none" },
          overflowX: "auto",
          scrollSnapType: "x mandatory",
          // display: "flex",
          gap: 1,
        }}
      >
        {!isLoading &&
          data?.blogs.map((item, i) => (
            <Box
              key={i}
              sx={{
                minWidth: "80%",
                height: 80,
                background: "#f0f0f0",
                scrollSnapAlign: "center",
                flexShrink: 0,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <Typography>{item.title ?? "Deleted"}</Typography>
              <Typography>{item.content ?? "Deleted"}</Typography>
            </Box>
          ))}
      </Box>
      <Pagination
        count={Math.ceil(data?.total ?? 0 / limit)}
        page={page}
        onChange={(_, value) => setPage(value)}
        size="small"
        sx={{ mt: 2 }}
      />
    </div>
  );
}

const tabs = ["Top All Time", "Top Trending"];
export default function RankingList() {
  const [currentTab, setCurrentTab] = useState<"allTime" | "trending">(
    "allTime"
  );
  const [currentTabIndex, setCurrentTabIndex] = useState(0);
  const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
    let action: string;
    switch (newValue) {
      case 0:
        action = "allTime";
        break;
      default:
        action = "trending";
        break;
    }
    setCurrentTabIndex(newValue);
    setCurrentTab(
      getTypeValidValue(action, ["allTime", "trending"], "allTime")
    );
  };
  return (
    <>
      <Box
        sx={{
          width: "100%",
        }}
      >
        <Tabs
          value={currentTabIndex}
          onChange={handleTabChange}
          indicatorColor="secondary"
          textColor="inherit"
          variant="fullWidth"
          aria-label="tabs"
        >
          {tabs.map((label, index) => (
            <Tab key={index} sx={{ fontWeight: "bold" }} label={label} />
          ))}
        </Tabs>
        <AllTimeBlogList type={currentTab} />
        <TrendingBlogList type={currentTab} />
      </Box>
    </>
  );
}
