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
  authors: Author[];
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

interface CommentNotification {
  notificationId: string;
  type: "authorNewContent";
  content: {
    authorID: string;
    authorName: string;
    authorSlug: string;
    urlSlug: string;
    title: string;
    content: string;
  };
  isRead?: boolean;
}

interface FollowNotification {
  notificationId: string;
  type: "follow";
  content: {
    followerId: string;
    followerName: string;
  };
  isRead?: boolean;
}

interface LikeNotification {
  notificationId: string;
  type: "like";
  content: {
    postId: string;
    likedBy: string;
  };
  isRead?: boolean;
}

export type Notification =
  | CommentNotification
  | FollowNotification
  | LikeNotification;

export interface Contact {
  contactID: string;
  email: string;
  content: string;
}
