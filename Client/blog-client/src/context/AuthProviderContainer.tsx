import { useQuery } from "@tanstack/react-query";
import { AuthProvider } from "./AuthContext";
import { fetchMe } from "../api/auth";
import { CircularProgress } from "@mui/material";
import { fetchAuthorMe } from "@/api/authorApi";

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
    staleTime: 1000 * 60 * 60 * 1,
  });

  const { data: author, isLoading: isLoadingAuthor } = useQuery({
    queryKey: ["me", "author"],
    queryFn: fetchAuthorMe,
    retry: false,
    enabled: !!user,
    refetchOnWindowFocus: false,
  });

  if (isLoading || isLoadingAuthor) {
    return <CircularProgress />;
  }

  return (
    <AuthProvider
      value={{
        user: user ?? null,
        isAuthenticated: !!user,
        authLoading: isLoading,
        author: author ?? null,
        isLoadingAuthor: isLoadingAuthor,
      }}
    >
      {children}
    </AuthProvider>
  );
}
