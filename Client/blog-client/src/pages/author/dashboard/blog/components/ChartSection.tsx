import { useMemo } from "react";
import Typography from "@mui/material/Typography";
import Grid from "@mui/material/Grid";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Slider from "@mui/material/Slider";
import Box from "@mui/material/Box";
import Chip from "@mui/material/Chip";
import { Line } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
  type ChartOptions,
} from "chart.js";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
  Legend,
);

interface ChartSectionProps<T extends Record<string, string | number>> {
  title: string;

  data: T[];

  periodField: keyof T;

  valueField: keyof T;

  comparisonLabel: string;

  sliderLabel: string;

  maxRange?: number;

  range: number;
  onRangeChange?: (range: number) => void;

  formatter?: (value: number) => string;
}

function createTrendLine(values: number[]): number[] {
  const n = values.length;

  if (n < 2) {
    return values;
  }

  const x = [...Array(n).keys()];

  const sumX = x.reduce((a, b) => a + b, 0);
  const sumY = values.reduce((a, b) => a + b, 0);

  const sumXY = x.reduce((acc, value, index) => acc + value * values[index], 0);

  const sumXX = x.reduce((acc, value) => acc + value * value, 0);

  const denominator = n * sumXX - sumX * sumX;

  if (denominator === 0) {
    return values;
  }

  const slope = (n * sumXY - sumX * sumY) / denominator;

  const intercept = (sumY - slope * sumX) / n;

  return x.map((value) => intercept + slope * value);
}

type MetricRecord = Record<string, string | number>;

export default function ChartSection<T extends MetricRecord>({
  title,
  data,
  periodField,
  valueField,
  comparisonLabel,
  sliderLabel,
  range,
  onRangeChange,
  maxRange = 7,
  formatter = (v) => v.toLocaleString(),
}: ChartSectionProps<T>) {
  if (!data || data.length == 0) {
    data = Array<T>(range);
  }
  const visibleData = useMemo(() => data.slice(-range), [data, range]);

  const labels = useMemo(
    () => visibleData.map((item) => String(item[periodField])),
    [visibleData, periodField],
  );

  const values = useMemo(
    () => visibleData.map((item) => Number(item[valueField])),
    [visibleData, valueField],
  );

  const trendLine = useMemo(() => createTrendLine(values), [values]);

  const current = values.at(-1) ?? 0;
  const previous = values.at(-2);

  const difference = previous == null ? 0 : current - previous;

  const percentChange =
    previous == null || previous === 0 ? 0 : (difference / previous) * 100;

  const isPositive = difference >= 0;

  const chartData = {
    labels,
    datasets: [
      {
        label: title,
        data: values,
        borderColor: "#1976d2",
        backgroundColor: "#1976d2",
        borderWidth: 2,
        tension: 0.35,
      },
      {
        label: "Trend",
        data: trendLine,
        borderColor: "#9e9e9e",
        borderDash: [6, 6],
        borderWidth: 2,
        pointRadius: 0,
      },
    ],
  };

  const options: ChartOptions<"line"> = {
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: "index",
    },
    plugins: {
      legend: {
        position: "top",
      },
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  const handleSliderChange = (_: Event, value: number | number[]) => {
    const nextRange = value as number;

    onRangeChange?.(nextRange);
  };

  return (
    <Box mb={4}>
      <Grid container spacing={2}>
        <Grid size={{ xs: 12, md: 9 }}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {title}
              </Typography>

              <Box height={320}>
                <Line data={chartData} options={options} />
              </Box>

              <Box mt={3}>
                <Typography variant="body2" gutterBottom>
                  Showing last {range}
                  {sliderLabel}
                </Typography>

                <Slider
                  min={1}
                  max={maxRange}
                  step={1}
                  value={range}
                  onChange={handleSliderChange}
                  valueLabelDisplay="auto"
                  marks={[
                    {
                      value: 1,
                      label: `1${sliderLabel}`,
                    },
                    {
                      value: Math.ceil(maxRange / 2),
                      label: `${Math.ceil(maxRange / 2)}${sliderLabel}`,
                    },
                    {
                      value: maxRange,
                      label: `${maxRange}${sliderLabel}`,
                    },
                  ]}
                />
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 12, md: 3 }}>
          <Card sx={{ height: "100%" }}>
            <CardContent>
              <Typography variant="overline" color="text.secondary">
                Trend
              </Typography>

              <Typography variant="body2" color="text.secondary" gutterBottom>
                {comparisonLabel}
              </Typography>

              <Typography variant="h4" sx={{ mt: 1 }}>
                {formatter(current)}
              </Typography>

              <Chip
                sx={{ mt: 2 }}
                color={isPositive ? "success" : "error"}
                label={`${isPositive ? "+" : ""}${percentChange.toFixed(1)}%`}
              />

              <Typography variant="body2" sx={{ mt: 2 }}>
                Change: {difference >= 0 ? "+" : ""}
                {formatter(difference)}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}
