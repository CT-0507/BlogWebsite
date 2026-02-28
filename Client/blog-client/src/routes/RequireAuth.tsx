import { useAuth } from "@/hooks/useAuth";
import { Navigate, Outlet, useLocation } from "react-router-dom";

type RequireAuthProps = {
  allowedRoles?: string[]; // allowed roles
};

function RequireAuth({ allowedRoles }: RequireAuthProps) {
  const location = useLocation();
  const { user, isAuthenticated, isLoading } = useAuth(); // replace with real logic

  // Still resolving auth
  if (isLoading) return <div>Loading...</div>;

  // Not logged in
  if (!isAuthenticated) {
    return <Navigate to="/account" replace state={{ from: location }} />;
  }

  // Role check
  if (allowedRoles && allowedRoles.length > 0) {
    const userRoles = Array.isArray(user?.roles)
      ? user.roles
      : user?.roles
      ? [user.roles]
      : [];

    const hasAccess = allowedRoles.some((r) => userRoles.includes(r));

    if (!hasAccess) {
      return <Navigate to="/403" replace />;
    }
  }

  // Authorized
  return <Outlet />;
}

export default RequireAuth;
