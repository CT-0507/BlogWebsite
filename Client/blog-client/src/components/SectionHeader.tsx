import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";

interface Props {
  title: string;
  subtitle?: string;
}

export default function SectionHeader({ title, subtitle }: Props) {
  return (
    <Stack spacing={1} mb={5}>
      <Typography variant="h3">{title}</Typography>

      {subtitle && <Typography color="text.secondary">{subtitle}</Typography>}
    </Stack>
  );
}
