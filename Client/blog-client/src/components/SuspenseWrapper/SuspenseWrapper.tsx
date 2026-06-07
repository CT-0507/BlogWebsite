import CircularProgress from "@mui/material/CircularProgress";
import { Suspense } from "react";

export default function SuspenseWrapper({ child }: { child: React.ReactNode }) {
  return <Suspense fallback={<CircularProgress />}>{child}</Suspense>;
}
