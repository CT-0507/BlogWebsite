import { useQuery } from "@tanstack/react-query";
import { AuthProvider } from "./AuthContext";
import { fetchMe } from "../api/auth";

export function AuthProviderContainer({
  children,
}: {
  children: React.ReactNode;
}) {
  const { data: user, isLoading } = useQuery({
    queryKey: ["me"],
    queryFn: fetchMe,
    retry: false,
  });

  return (
    <AuthProvider
      value={{
        user: user ?? null,
        isAuthenticated: !!user,
        isLoading,
      }}
    >
      {children}
    </AuthProvider>
  );
}
