import { Routes, Route } from "react-router-dom";
import Home from "@/pages/home/Home";
import BasicLayout from "@/layouts/BasicLayout";
import RequireAuth from "./RequireAuth";
import Dashboard from "@/pages/dashboard/Dashboard";
import { ROLES } from "@/config/roles";
import UserHome from "@/pages/user/UserHome";
import Profile from "@/pages/user/profile/Profile";
import PublishPage from "@/pages/author/dashboard/blog/publish/Publish";
import NotFound from "@/pages/NotFound/NotFound";
import ErrorBoundary from "@/pages/error/ErrorBoundary";
import ViewBlog from "@/pages/blog/viewBlog/ViewBlogPage";
import AuthorBlog from "@/pages/blog/authorBlog/AuthorBlog";
import CreateAuthorPage from "@/pages/author/create-author/CreateAuthorPage";
import ResponsiveSidebarLayout from "@/pages/author/dashboard/layout";
import SuspenseWrapper from "@/components/SuspenseWrapper/SuspenseWrapper";
import { lazy } from "react";
import Unauthorized from "@/pages/Unauthorized/Unauthorized";

const AuthPage = lazy(() => import("@/pages/auth/Auth"));
const AuthorDashboard = lazy(
  () => import("@/pages/author/dashboard/Dashboard"),
);
const MyBlogs = lazy(() => import("@/pages/author/dashboard/blogs/MyBlogs"));
const EditPage = lazy(
  () => import("@/pages/author/dashboard/blog/edit/EditBlog"),
);
const ViewBlogDashboardPage = lazy(
  () => import("@/pages/author/dashboard/blog/view/ViewBlog"),
);
const PortfolioPage = lazy(
  () => import("@/pages/about/portfolio/PortfolioPage"),
);
const AboutPage = lazy(() => import("@/pages/about/AboutPage"));

export default function AppRoutes() {
  return (
    <ErrorBoundary>
      <Routes>
        <Route path="/" element={<BasicLayout />}>
          <Route index element={<Home />} />
          <Route
            path="account"
            element={<SuspenseWrapper child={<AuthPage />} />}
          />
          <Route
            path="author"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
          >
            <Route path="" element={<ResponsiveSidebarLayout />}>
              <Route
                path="dashboard"
                element={<SuspenseWrapper child={<AuthorDashboard />} />}
              />
              <Route
                path="my-blogs"
                element={<SuspenseWrapper child={<MyBlogs />} />}
              />
              <Route
                path="my-blogs/:slug"
                element={<SuspenseWrapper child={<ViewBlogDashboardPage />} />}
              />
              <Route
                path="my-blogs/:slug/edit"
                element={<SuspenseWrapper child={<EditPage />} />}
              />
              <Route
                path="my-blogs/:slug/view"
                element={<SuspenseWrapper child={<ViewBlogDashboardPage />} />}
              />

              <Route path="401" element={<Unauthorized />} />
              <Route path="*" element={<NotFound />} />
            </Route>
          </Route>
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
            path="authors"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
          >
            <Route path="new" element={<CreateAuthorPage />} />
          </Route>
          <Route
            path="blogs"
            element={<RequireAuth allowedRoles={[ROLES.ADMIN, ROLES.USER]} />}
          >
            <Route path="publish" element={<PublishPage />} />
          </Route>
          <Route path="blogs">
            <Route index element={<Home />} />
            <Route path=":slug" element={<ViewBlog />} />
            <Route path="author">
              <Route path=":slug" element={<AuthorBlog />} />
            </Route>
          </Route>
          <Route path="about">
            <Route index element={<SuspenseWrapper child={<AboutPage />} />} />
            <Route
              path="creator"
              element={<SuspenseWrapper child={<PortfolioPage />} />}
            />
          </Route>
          {/* Fallback */}
          <Route path="/401" element={<Unauthorized />} />
          <Route path="*" element={<NotFound />} />
        </Route>
      </Routes>
    </ErrorBoundary>
  );
}
