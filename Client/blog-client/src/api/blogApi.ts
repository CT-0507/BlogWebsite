import type { PublishBlogFormValues } from "@/pages/author/dashboard/blog/publish/model/schema";
import { api, API_VERSION_V1, axiosAuth } from "./axiosConfig";
import type { PostCommentFormValues } from "@/pages/blog/viewBlog/model/schema";
import type {
  Blog,
  BlogComment,
  BlogReaction,
  BlogReactionType,
  CommentReaction,
  RankingBlogData,
} from "@/types/Blog";
import { getQueryParam } from "@/utils/mapper";
import type { BlogReport } from "@/types/types";

export async function publishBlogRequest(
  formData: PublishBlogFormValues & {
    files: Map<string, File>;
    idempotencyKey: string;
  },
): Promise<Blog> {
  const formDataV = new FormData();
  formDataV.append("title", formData.title);
  formDataV.append("urlSlug", formData.urlSlug);
  formDataV.append("contentText", formData.content.plainText);
  formDataV.append("contentJson", JSON.stringify(formData.content.json));

  formData.files.forEach((file, tempId) => {
    formDataV.append(tempId, file);
  });

  if (formData.thumbnail) {
    formDataV.append("thumbnail", formData.thumbnail);
  }
  if (formData.tags) {
    formData.tags.forEach((item) => {
      formDataV.append("tags", item);
    });
  }
  const { data } = await axiosAuth.post(`${API_VERSION_V1}/blogs`, formDataV, {
    headers: {
      "Idempotency-Key": formData.idempotencyKey,
    },
  });

  return data;
}

export type UpdateBlogRequestParams = PublishBlogFormValues & {
  files: Map<string, File>;
  blogId: number;
  idempotencyKey: string;
};
export async function updateBlogRequest(
  formData: UpdateBlogRequestParams,
): Promise<Blog> {
  const formDataV = new FormData();
  formDataV.append("title", formData.title);
  formDataV.append("urlSlug", formData.urlSlug);
  formDataV.append("contentText", formData.content.plainText);
  formDataV.append("contentJson", JSON.stringify(formData.content.json));

  formData.files.forEach((file, tempId) => {
    formDataV.append(tempId, file);
  });

  if (formData.thumbnail) {
    formDataV.append("thumbnail", formData.thumbnail);
  }
  if (formData.tags) {
    formData.tags.forEach((item) => {
      formDataV.append("tags", item);
    });
  }
  const { data } = await axiosAuth.patch(
    `${API_VERSION_V1}/blogs/${formData.blogId}`,
    formDataV,
  );

  return data;
}

interface QueryBlogsParams {
  title?: string | null;
  content?: string | null;
  author?: string | null;
  sortBy: string;
  sortDir: string;
  limit: number;
}

interface ListBlogsResponse {
  total: number;
  blogs: Blog[];
}

export async function listBlogs(
  queryParams: QueryBlogsParams,
  page: number,
): Promise<ListBlogsResponse> {
  const params = getQueryParam(queryParams);

  params.append("page", page.toString());

  const { data } = await api.get(
    `${API_VERSION_V1}/blogs?` + params.toString(),
  );

  return data;
}

interface ListLikedBlogsResponse {
  total: number;
  blogs: Blog[];
}

export async function listLikedBlogs(): Promise<ListLikedBlogsResponse> {
  const { data } = await axiosAuth.get(
    `${API_VERSION_V1}/blogs/me/liked-blogs`,
  );

  return data;
}

export async function listMyBlogs(
  queryParams: Omit<QueryBlogsParams, "author">,
  page: number,
): Promise<ListBlogsResponse> {
  const params = getQueryParam(queryParams);

  params.append("page", page.toString());

  const { data } = await axiosAuth.get(
    `${API_VERSION_V1}/dashboard/author/blogs?` + params.toString(),
  );

  return data;
}

export async function listBlogsAuthor(
  queryParams: Omit<QueryBlogsParams, "author">,
  page: number,
): Promise<ListBlogsResponse> {
  const params = getQueryParam(queryParams);

  params.append("page", page.toString());

  const { data } = await axiosAuth.get(
    `${API_VERSION_V1}/dashboard/blogs?` + params.toString(),
  );

  return data;
}

export async function getBlogBySlug(
  slug: string,
  isAuthenticated?: boolean,
): Promise<Blog> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(
      `${API_VERSION_V1}/blogs/slug/${slug}`,
    );

    return data;
  }

  const { data } = await api.get(`${API_VERSION_V1}/blogs/slug/${slug}`);

  return data;
}

export interface PostCommentParams {
  formData: PostCommentFormValues;
  blogID: string;
}

// Comments API
export async function postComment(
  formData: PostCommentFormValues,
): Promise<BlogComment> {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/blogs/${formData.blogID}/comments`,
    formData,
  );

  return data;
}

export interface GetRootCommentsResponse {
  total: number;
  comments: BlogComment[];
}

export async function getRootComments(
  blogID: number,
  isAuthenticated?: boolean,
): Promise<GetRootCommentsResponse> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(
      `${API_VERSION_V1}/blogs/${blogID}/comments`,
    );

    return data;
  }

  const { data } = await api.get(`${API_VERSION_V1}/blogs/${blogID}/comments`);

  return data;
}

export async function getReplies(
  parentID: string,
  isAuthenticated?: boolean,
): Promise<BlogComment[]> {
  if (isAuthenticated) {
    const { data } = await axiosAuth.get(
      `${API_VERSION_V1}/comments/${parentID}/children`,
    );

    return data;
  }

  const { data } = await api.get(
    `${API_VERSION_V1}/comments/${parentID}/children`,
  );

  return data;
}

export interface CreateBlogReactionResponse {
  transitionType: string;
  blogId: number;
  type: BlogReactionType;
}

export async function createBlogReaction(
  formData: BlogReaction,
): Promise<CreateBlogReactionResponse> {
  const { data } = await axiosAuth.post(
    `${API_VERSION_V1}/blogs/${formData.blogId}/reaction`,
    {
      type: formData.type,
    },
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
    `${API_VERSION_V1}/comments/${formData.commentId}/reaction`,
    {
      type: formData.type,
    },
  );

  return data;
}

export interface UpdateBlogCommentContentRequest {
  content: string;
  commentId: string;
}

interface UpdateCommentContentResponse {
  commentId: string;
  content: string;
}

export async function updateCommentContent(
  formData: UpdateBlogCommentContentRequest,
): Promise<UpdateCommentContentResponse> {
  const { data } = await axiosAuth.patch(
    `${API_VERSION_V1}/comments/${formData.commentId}`,
    {
      content: formData.content,
    },
  );

  return data;
}

export async function hideComment(commentId: string) {
  const { data } = await axiosAuth.patch(
    `${API_VERSION_V1}/comments/${commentId}/hidden`,
  );

  return data;
}

export async function deleteComment(commentId: string) {
  const { data } = await axiosAuth.delete(
    `${API_VERSION_V1}/comments/${commentId}/delete`,
  );

  return data;
}

interface GetTrendingBlogsParams {
  limit?: number;
  sortBy?: string;
  sortDir?: string;
}

interface GetTrendingBlogsResponse {
  total: number;
  blogs: RankingBlogData[];
}

export async function getRankingBlogs(
  queryParams: GetTrendingBlogsParams,
  page: number,
  type: "allTime" | "trending",
): Promise<GetTrendingBlogsResponse> {
  const params = getQueryParam(queryParams);

  params.append("page", page.toString());

  const { data } = await api.get(
    `${API_VERSION_V1}/blogs/ranking?type=${type}&` + params.toString(),
  );

  return data;
}

export async function uploadByFile(file: File) {
  const formData = new FormData();
  formData.append("image", file);

  const res = await axiosAuth.post(`${API_VERSION_V1}/uploads/image`, formData);

  return res.data;
}

interface CreateBlogReportRequest {
  blogID: number;
  reason: string;
}

export async function createBlogReport(
  report: CreateBlogReportRequest,
): Promise<BlogReport> {
  const res = await axiosAuth.post(
    `${API_VERSION_V1}/blogs/${report.blogID}/reports`,
    {
      reason: report.reason,
    },
  );

  return res.data;
}

interface GetBlogMetricsRequest {
  resultLength?: number;
  viewType: string;
  blogID: number;
}

export async function getBlogMetrics({
  blogID,
  viewType,
  resultLength = 4,
}: GetBlogMetricsRequest) {
  const params = getQueryParam({ viewType, resultLength });
  const res = await axiosAuth.get(
    `${API_VERSION_V1}/dashboard/author/blogs/${blogID}/metrics?` +
      params.toString(),
  );

  return res.data;
}
