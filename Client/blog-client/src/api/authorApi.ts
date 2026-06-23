import type { Author } from "@/types/types";
import type { CreateAuthorFormValues } from "./../pages/author/create-author/model/schema";
import { api, axiosAuth, API_VERSION_V1 } from "./axiosConfig";

export async function getAuthorProfileRequest(slug: string) {
  const { data } = await api.get(`${API_VERSION_V1}/authors/` + slug);

  return data;
}

export async function getAuthorBlogsRequest(slug: string) {
  const { data } = await api.get(`${API_VERSION_V1}/blogs/author/slug/` + slug);

  return data;
}

interface GetFollowedAuthorsResponse {
  length: number;
  authors: Author[];
}

export async function getFollowedAuthorsRequest(): Promise<GetFollowedAuthorsResponse> {
  const { data } = await axiosAuth.get(
    `${API_VERSION_V1}/authors/me/following/authors`,
  );

  return data;
}

export async function getAuthorFollowersRequest(slug: string) {
  const { data } = await api.get(`${API_VERSION_V1}/authors/${slug}/followers`);

  return data;
}

export async function createAuthorRequest(formData: CreateAuthorFormValues) {
  const formDataV = new FormData();
  formDataV.append("displayName", formData.displayName);
  if (formData.bio) formDataV.append("bio", formData.bio);
  formDataV.append("slug", formData.slug);
  if (formData.socialLink) formDataV.append("socialLink", formData.socialLink);
  if (formData.email) formDataV.append("email", formData.email);
  if (formData.avatar) formDataV.append("avatar", formData.avatar);
  const { data } = await axiosAuth.post(`${API_VERSION_V1}/authors`, formDataV);

  return data;
}

export async function followAuthorRequest(authorID: string) {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/authors/${authorID}/follow`,
  );

  return data;
}

export async function unfollowAuthorRequest(authorID: string) {
  const { data } = await axiosAuth.delete(
    `${API_VERSION_V1}/authors/${authorID}/follow`,
  );

  return data;
}

export async function fetchAuthorMe(): Promise<Author> {
  const { data } = await axiosAuth.get(`${API_VERSION_V1}/me/authorProfile`);
  return data || null;
}

type DashboardViewMetrics = {
  todayViews: number;
  yesterdayViews: number;
  thisWeekViews: number;
  lastWeekViews: number;
};

type DashboardReactionMetrics = {
  todayLikes: number;
  todayDislikes: number;
  yesterdayLikes: number;
  yesterdayDislikes: number;
  thisWeekLikes: number;
  thisWeekDislikes: number;
  lastWeekLikes: number;
  lastWeekDislikes: number;
};

export type AuthorDashboardMetrics = {
  viewsMetrics: DashboardViewMetrics;
  reactionMetrics: DashboardReactionMetrics;
};

export async function getAuthorDashboardMetrics(): Promise<AuthorDashboardMetrics> {
  const { data } = await axiosAuth.get(`${API_VERSION_V1}/dashboard/author`);
  return data || null;
}
