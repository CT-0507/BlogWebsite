import { Routes, Route } from "react-router-dom";
import Home from "@/pages/home/Home";
import Auth from "@/pages/auth/Auth";
import BasicLayout from "@/layouts/BasicLayout";
import RequireAuth from "./RequireAuth";
import Dashboard from "@/pages/dashboard/Dashboard";
import { ROLES } from "@/config/roles";
import UserHome from "@/pages/user/UserHome";
import Profile from "@/pages/user/profile/Profile";
import PublishPage from "@/pages/blog/publish/Publish";
import NotFound from "@/pages/NotFound/NotFound";
import ErrorBoundary from "@/pages/error/ErrorBoundary";
import ViewBlog from "@/pages/blog/viewBlog/ViewBlog";
import AuthorBlog from "@/pages/blog/authorBlog/authorBlog";

export default function AppRoutes() {
  return (
    <ErrorBoundary>
      <Routes>
        <Route path="/" element={<BasicLayout />}>
          <Route index element={<Home />} />
          <Route path="account" element={<Auth />} />
          <Route
            path="dashboard"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN]} />}
          >
            <Route index element={<Dashboard />} />
          </Route>
          <Route
            path="user"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
          >
            <Route index element={<UserHome />} />
            <Route path="profile" element={<Profile />} />
          </Route>
          <Route
            path="blogs"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
          >
            <Route path="publish" element={<PublishPage />} />
          </Route>
          <Route path="blogs">
            <Route index element={<Home />} />
            <Route path=":id" element={<ViewBlog />} />
            <Route path="author">
              <Route path=":id" element={<AuthorBlog />} />
            </Route>
          </Route>
        </Route>
        {/* Fallback */}
        <Route path="*" element={<NotFound />} />
      </Routes>
    </ErrorBoundary>
  );
}
