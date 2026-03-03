import { Routes, Route } from "react-router-dom";
import Home from "@/pages/home/Home";
import Auth from "@/pages/auth/Auth";
import BasicLayout from "@/layouts/BasicLayout";
import RequireAuth from "./RequireAuth";
import Dashboard from "@/pages/dashboard/Dashboard";
import { ROLES } from "@/config/roles";
import UserHome from "@/pages/user/UserHome";
import Profile from "@/pages/user/profile/Profile";

export default function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<BasicLayout />}>
        <Route index element={<Home />} />
        <Route path="/account" element={<Auth />} />
        <Route
          path="/dashboard"
          element={<RequireAuth allowedRoles={[ROLES.ADMIN]} />}
        >
          <Route index element={<Dashboard />} />
        </Route>
        <Route
          path="/user"
          element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
        >
          <Route index element={<UserHome />} />
          <Route path="profile" element={<Profile />} />
        </Route>
      </Route>
    </Routes>
  );
}
