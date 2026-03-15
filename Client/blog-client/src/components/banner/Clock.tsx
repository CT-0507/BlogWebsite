import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { useEffect, useState } from "react";

export function ClockBanner() {
  const [time, setTime] = useState(new Date());

  useEffect(() => {
    const id = setInterval(() => setTime(new Date()), 1000);
    return () => clearInterval(id);
  }, []);

  return (
    <Box sx={{ textAlign: "center", py: 1 }}>
      <Typography variant="caption" color="text.secondary">
        Current time: {time.toLocaleTimeString()}
      </Typography>
    </Box>
  );
}
