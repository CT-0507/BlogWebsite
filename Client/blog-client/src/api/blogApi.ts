import type { PublishBlogFormValues } from "@/pages/blog/publish/model/schema";
import { api, axiosAuth } from "./axiosConfig";

export async function publishBlogRequest(formData: PublishBlogFormValues) {
  const { data } = await axiosAuth.post("/blogs", formData);

  return data;
}

export async function listBlogs(queryParams: string) {
  const { data } = await api.get("/blogs" + queryParams);

  return data;
}

export async function getBlogByID(id: string) {
  const { data } = await api.get(`/blogs/${id}`);

  return data;
}
