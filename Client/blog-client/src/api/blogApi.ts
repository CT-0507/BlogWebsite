import type { PublishBlogFormValues } from "@/pages/blog/publish/model/schema";
import { axiosAuth } from "./axiosConfig";

export async function publishBlogRequest(formData: PublishBlogFormValues) {
  const { data } = await axiosAuth.post("/blogs", formData);

  return data;
}
