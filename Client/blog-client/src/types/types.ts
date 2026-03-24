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

export type { Author, FollowedAuthorReponse };
