import type { PublishBlogFormValues } from "@/pages/blog/publish/model/schema";
import { api, axiosAuth } from "./axiosConfig";
import { isLocalMode } from ".";
import { blogPOST, blogs } from "./mockApi";

export async function publishBlogRequest(formData: PublishBlogFormValues) {
  if (isLocalMode) return blogPOST;

  const { data } = await axiosAuth.post("/blogs", formData);

  return data;
}

export async function listBlogs(queryParams: string) {
  if (isLocalMode) return blogs;
  const { data } = await api.get("/blogs" + queryParams);

  return data;
}

export async function getBlogBySlug(slug: string) {
  if (isLocalMode) return blogs.find((i) => i.urlSlug == slug);
  const { data } = await api.get(`/blogs/${slug}`);

  return data;
}
