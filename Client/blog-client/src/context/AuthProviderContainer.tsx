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
  });

  const { data: author, isLoading: isLoadingAuthor } = useQuery({
    queryKey: ["me", "author"],
    queryFn: fetchAuthorMe,
    retry: false,
    enabled: !!user,
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
