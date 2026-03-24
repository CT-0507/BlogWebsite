import type { CreateAuthorFormValues } from "./../pages/author/model/schema";
import { api, axiosAuth } from "./axiosConfig";

const API_VERSION = "/api/v1";

export async function getAuthorProfileRequest(slug: string) {
  const { data } = await api.get(`${API_VERSION}/authors/` + slug);

  return data;
}

export async function getAuthorBlogsRequest(slug: string) {
  const { data } = await api.get("/blogs/author/" + slug);

  return data;
}

export async function getFollowedAuthorsRequest() {
  const { data } = await axiosAuth.get(
    `${API_VERSION}/authors/me/following/authors`
  );

  return data;
}

export async function getAuthorFollowersRequest(slug: string) {
  const { data } = await api.get(`${API_VERSION}/authors/${slug}/followers`);

  return data;
}

export async function createAuthorRequest(formData: CreateAuthorFormValues) {
  const { data } = await axiosAuth.post(`${API_VERSION}/authors`, formData);

  return data;
}

export async function followAuthorRequest(authorID: string) {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/authors/${authorID}/follow`
  );

  return data;
}

export async function unfollowAuthorRequest(authorID: string) {
  const { data } = await axiosAuth.delete(
    `${API_VERSION}/authors/${authorID}/follow`
  );

  return data;
}
