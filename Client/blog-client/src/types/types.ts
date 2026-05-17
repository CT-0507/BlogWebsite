interface Author {
  authorID: string;
  displayName: string;
  bio: string;
  avatar: string;
  slug: string;
  socialLink: string;
  email: string;
  followerCount: number;
  blogCount: number;
  createdAt: string;
}

interface FollowedAuthorReponse {
  message: string;
  length: number;
  authors: string[];
}

export const BLOG_SORT_BY_VALUES = [
  "title",
  "created_at",
  "relevance",
] as const;

export const SORT_DIR = ["asc", "desc"] as const;

export type SortByValue = (typeof BLOG_SORT_BY_VALUES)[number];
export type SortDir = (typeof SORT_DIR)[number];

export type { Author, FollowedAuthorReponse };

export const BLOG_RANKING_SORT_BY_VALUES = [
  "daily",
  "weekly",
  "likes",
  "score",
  "rank",
] as const;
export type SortBlogRankingByValue =
  (typeof BLOG_RANKING_SORT_BY_VALUES)[number];

export type BlogReport = {
  reportID: number;
  blogID: number;
  reason: string;
};
