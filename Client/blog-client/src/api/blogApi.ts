import type { PublishBlogFormValues } from "@/pages/blog/publish/model/schema";
import { api, axiosAuth } from "./axiosConfig";
import type { PostCommentFormValues } from "@/pages/blog/viewBlog/model/schema";
import type {
  Blog,
  BlogComment,
  BlogReaction,
  BlogReactionType,
  CommentReaction,
} from "@/types/Blog";

const API_VERSION = "/api/v1";

export async function publishBlogRequest(formData: PublishBlogFormValues) {
  const { data } = await axiosAuth.post(`${API_VERSION}/blogs`, formData);

  return data;
}

export async function listBlogs(queryParams: string) {
  const { data } = await api.get(`${API_VERSION}/blogs` + queryParams);

  return data;
}

export async function getBlogBySlug(
  slug: string,
  isAuthenticated?: boolean
): Promise<Blog> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(`${API_VERSION}/blogs/slug/${slug}`);

    return data;
  }

  const { data } = await api.get(`${API_VERSION}/blogs/slug/${slug}`);

  return data;
}

export interface PostCommentParams {
  formData: PostCommentFormValues;
  blogID: string;
}

// Comments API
export async function postComment(
  formData: PostCommentFormValues
): Promise<BlogComment> {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/blogs/${formData.blogID}/comments`,
    formData
  );

  return data;
}

export interface GetRootCommentsResponse {
  total: number;
  comments: BlogComment[];
}

export async function getRootComments(
  blogID: number,
  isAuthenticated?: boolean
): Promise<GetRootCommentsResponse> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(
      `${API_VERSION}/blogs/${blogID}/comments`
    );

    return data;
  }

  const { data } = await api.get(`${API_VERSION}/blogs/${blogID}/comments`);

  return data;
}

export async function getReplies(
  parentID: string,
  isAuthenticated?: boolean
): Promise<BlogComment[]> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(
      `${API_VERSION}/comments/${parentID}/children`
    );

    return data;
  }

  const { data } = await api.get(
    `${API_VERSION}/comments/${parentID}/children`
  );

  return data;
}

export interface CreateBlogReactionResponse {
  transitionType: string;
  blogId: number;
  type: BlogReactionType;
}

export async function createBlogReaction(
  formData: BlogReaction
): Promise<CreateBlogReactionResponse> {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/blogs/${formData.blogId}/reaction`,
    {
      type: formData.type,
    }
  );

  return data;
}

export interface CreateCommentReactionResponse {
  transitionType: string;
  commentId: string;
  type: BlogReactionType;
}

export async function createCommentReaction(formData: CommentReaction) {
  const { data } = await axiosAuth.post(
    `${API_VERSION}/comments/${formData.commentId}/reaction`,
    {
      type: formData.type,
    }
  );

  return data;
}
