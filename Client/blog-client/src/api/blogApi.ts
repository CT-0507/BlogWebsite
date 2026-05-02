import type { PublishBlogFormValues } from "@/pages/blog/publish/model/schema";
import { api, axiosAuth } from "./axiosConfig";
import type { PostCommentFormValues } from "@/pages/blog/viewBlog/model/schema";
import type { Blog } from "@/pages/home/BlogList";
import type { BlogComment } from "@/pages/blog/viewBlog/CommentSection";

const API_VERSION = "/api/v1";

export async function publishBlogRequest(formData: PublishBlogFormValues) {
  const { data } = await axiosAuth.post(`${API_VERSION}/blogs`, formData);

  return data;
}

export async function listBlogs(queryParams: string) {
  const { data } = await api.get(`${API_VERSION}/blogs` + queryParams);

  return data;
}

export async function getBlogBySlug(slug: string) {
  const { data } = await api.get(`${API_VERSION}/blogs/slug/${slug}`);

  return data as Blog;
}

export interface PostCommentParams {
  formData: PostCommentFormValues;
  blogID: string;
}

// Comments API
export async function postComment(formData: PostCommentFormValues) {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/blogs/${formData.blogID}/comments`,
    formData
  );

  return data;
}

export async function getRootComments(blogID: number): Promise<BlogComment[]> {
  const { data } = await api.get(`${API_VERSION}/blogs/${blogID}/comments`);

  return data;
}

export async function getReplies(parentID: string): Promise<BlogComment[]> {
  const { data } = await api.get(
    `${API_VERSION}/comments/${parentID}/children`
  );

  return data;
}

export async function likeBlog(blogID: string) {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/blogs/${blogID}/reaction`
  );

  return data;
}

export async function dislikeBlog(blogID: string) {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/blogs/${blogID}/reaction`
  );

  return data;
}
