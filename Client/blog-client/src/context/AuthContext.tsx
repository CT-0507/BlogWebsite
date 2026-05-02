import { createContext } from "react";

export type User = {
  userID: string;
  username: string;
  roles: string[];
  email: string;
  firstName: string;
  lastName: string;
};

type AuthContextType = {
  user: User | null;
  isAuthenticated: boolean;
  authLoading: boolean;
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({
  children,
  value,
}: {
  children: React.ReactNode;
  value: AuthContextType;
}) {
  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export default AuthContext;
