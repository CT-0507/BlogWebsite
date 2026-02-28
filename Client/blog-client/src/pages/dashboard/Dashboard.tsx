import { axiosAuth } from "@/api/axiosConfig";
import Typography from "@mui/material/Typography";
import { useQuery } from "@tanstack/react-query";

export default function Dashboard() {
  const { data } = useQuery({
    queryKey: ["nothing"],
    queryFn: async () => {
      const { data } = await axiosAuth.get("/dashboard");
      return data;
    },
  });
  return (
    <>
      <Typography component="h1">{data?.message}</Typography>
    </>
  );
}
