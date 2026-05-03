export interface Blog {
  blogID: number;
  author: {
    authorID: string;
    slug: string;
    displayName: string;
  };
  title: string;
  likeCount: number;
  dislikeCount: number;
  urlSlug: string;
  content: string;
  userReaction?: string | null;
  createdAt: string;
}

export interface BlogComment {
  commentId: string;
  blogId: number;
  actorId: string;
  actorType: string;
  actorAvatarUrl?: string | null;
  actorDisplayName: string;
  content: string;
  userReaction?: string | null;
  likeCount: number;
  dislikeCount: number;
  replyCount: number;
  rootCommentId: string;
  createdAt: string;
  updatedAt: string;
  status?: string | null;
  parentCommentId?: string | null;
}

export type BlogReactionType = "like" | "dislike";

export interface BlogReaction {
  blogId: number;
  userId: string;
  type: BlogReactionType;
}

export type CommentReactionType = "like" | "dislike";

export interface CommentReaction {
  commentId: string;
  userId: string;
  type: CommentReactionType;
}
