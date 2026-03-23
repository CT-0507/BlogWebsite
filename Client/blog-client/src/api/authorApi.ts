import { axiosAuth } from "./axiosConfig";

export async function getAuthorProfile(slug: string) {
  const { data } = await axiosAuth.get("/api/v1/authors/" + slug);

  return data;
}

export async function getAuthorBlogs(slug: string) {
  const { data } = await axiosAuth.get("/blogs/author/" + slug);

  return data;
}
