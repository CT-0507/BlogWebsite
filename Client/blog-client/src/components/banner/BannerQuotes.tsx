import Typography from "@mui/material/Typography";
import Box from "@mui/material/Box";
import { useMemo, useState } from "react";

const quotes = [
  {
    text: "Simplicity is the ultimate sophistication.",
    author: "Leonardo da Vinci",
  },
  {
    text: "What we think, we become.",
    author: "Buddha",
  },
  {
    text: "Stay hungry, stay foolish.",
    author: "Steve Jobs",
  },
  {
    text: "The only way to do great work is to love what you do.",
    author: "Steve Jobs",
  },
  {
    text: "In the middle of difficulty lies opportunity.",
    author: "Albert Einstein",
  },
];

export default function QuoteBanner() {
  const [randomNumber] = useState(() =>
    Math.floor(Math.random() * quotes.length)
  );
  const quote = useMemo(() => {
    return quotes[Math.floor(randomNumber)];
  }, [randomNumber]);

  return (
    <Box
      sx={{
        width: "100%",
        bgcolor: "info.main",
        color: "primary.contrastText",
        py: { xs: 4, md: 6 },
        px: 2,
        textAlign: "center",
      }}
    >
      <Typography
        variant="h6"
        sx={{ fontStyle: "italic", maxWidth: 900, mx: "auto" }}
      >
        “{quote?.text}”
      </Typography>

      <Typography variant="body2" sx={{ mt: 1, opacity: 0.8 }}>
        — {quote?.author}
      </Typography>
    </Box>
  );
}
