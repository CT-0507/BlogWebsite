import { Link as RouterLink, useLocation } from "react-router-dom";
import Breadcrumbs from "@mui/material/Breadcrumbs";
import Link from "@mui/material/Link";
import Typography from "@mui/material/Typography";

function formatLabel(segment: string) {
  return segment
    .replace(/-/g, " ")
    .replace(/_/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

interface NavBreadcrumbsProps {
  hideLast?: number;
  hiddenSegments?: string[];
}

export default function NavBreadcrumbs({
  hideLast = 0,
  hiddenSegments = [],
}: NavBreadcrumbsProps) {
  const location = useLocation();

  const allSegments = location.pathname
    .split("/")
    .filter(Boolean)
    .filter((segment) => {
      // Hide actions
      if (hiddenSegments.includes(segment)) return false;

      // Hide slug-like segment when preceded by my-blogs
      const idx = location.pathname.split("/").filter(Boolean).indexOf(segment);

      const previous = location.pathname.split("/").filter(Boolean)[idx - 1];

      if (previous === "my-blogs") return false;

      return true;
    });

  const visibleSegments =
    hideLast > 0 ? allSegments.slice(0, -hideLast) : allSegments;

  return (
    <Breadcrumbs aria-label="breadcrumb">
      <Link component={RouterLink} underline="hover" color="inherit" to="/">
        Home
      </Link>

      {visibleSegments.map((segment, index) => {
        const to = `/${allSegments.slice(0, index + 1).join("/")}`;

        const isLast = index === visibleSegments.length - 1;

        return isLast ? (
          <Typography key={to} color="text.primary">
            {formatLabel(segment)}
          </Typography>
        ) : (
          <Link
            key={to}
            component={RouterLink}
            underline="hover"
            color="inherit"
            to={to}
          >
            {formatLabel(segment)}
          </Link>
        );
      })}
    </Breadcrumbs>
  );
}
