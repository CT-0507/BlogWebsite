import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Typography from "@mui/material/Typography";
import LinearProgress from "@mui/material/LinearProgress";
import Stack from "@mui/material/Stack";

interface Props {
  name: string;
  value: number;
}

export default function SkillProgress({ name, value }: Props) {
  return (
    <Card>
      <CardContent>
        <Stack spacing={2}>
          <Typography>{name}</Typography>

          <LinearProgress variant="determinate" value={value} />

          <Typography variant="body2">{value}%</Typography>
        </Stack>
      </CardContent>
    </Card>
  );
}
