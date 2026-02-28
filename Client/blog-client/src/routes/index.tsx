import { Routes, Route } from "react-router-dom";
import Home from "@/pages/home/Home";
import Auth from "@/pages/auth/Auth";
import BasicLayout from "@/layouts/BasicLayout";
import RequireAuth from "./RequireAuth";
import Dashboard from "@/pages/dashboard/Dashboard";
import { ROLES } from "@/config/roles";

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
      </Route>
    </Routes>
  );
}
