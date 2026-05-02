import { useQuery } from "@tanstack/react-query";
import { AuthProvider } from "./AuthContext";
import { fetchMe } from "../api/auth";
import { CircularProgress } from "@mui/material";

export function AuthProviderContainer({
  children,
}: {
  children: React.ReactNode;
}) {
  const hasSession = localStorage.getItem("hasSession") === "true";

  const { data: user, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: fetchMe,
    retry: false,
    enabled: hasSession,
  });

  if (isLoading) {
    return <CircularProgress />;
  }

  return (
    <AuthProvider
      value={{
        user: user ?? null,
        isAuthenticated: !!user,
        authLoading: isLoading,
      }}
    >
      {children}
    </AuthProvider>
  );
}
