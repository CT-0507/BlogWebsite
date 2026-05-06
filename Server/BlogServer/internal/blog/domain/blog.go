package domain

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

type Blog struct {
	BlogID   int64  `json:"id"`
	AuthorID string `json:"authorID"`
	Title    string `json:"title"`
	URLSlug  string `json:"urlSlug"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	// Tags    []Tag  `json:"tags"`
	// Images  []string
	model.Audit
}

type BlogWithAuthorData struct {
	BlogID       int64      `json:"blogID"`
	Author       AuthorData `json:"author"`
	URLSlug      string     `json:"urlSlug"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	LikeCount    int64      `json:"likeCount"`
	DislikeCount int64      `json:"dislikeCount"`
	UserReaction *string    `json:"userReaction"`
	Status       string     `json:"status"`
	model.Audit
}

type AuthorData struct {
	AuthorID    string  `json:"authorID"`
	AvatarURL   *string `json:"avatarURL"`
	Slug        string  `json:"slug"`
	DisplayName string  `json:"displayName"`
}

type RankingBlogData struct {
	BlogID              int64    `json:"blogID"`
	TotalAllTimeResult  *int64   `json:"totalAllTime,omitempty"`
	TotalTrendingResult *int64   `json:"totalTrending,omitempty"`
	RankAllTime         *int32   `json:"rankAllTime"`
	RankTrending        *int32   `json:"rankTrending"`
	ScoreAllTime        *float64 `json:"scoreAllTime"`
	ScoreTrending       *float64 `json:"scoreTrending"`
	LikeCount           int32    `json:"likeCount"`
	DislikeCount        int32    `json:"dislikeCount"`
	CommentCount        int32    `json:"commentCount"`
	WeeklyAccessCount   int32    `json:"weeklyAccessCount"`
	DailyAccessCount    int32    `json:"dailyAccessCount"`
	CreatedAt           string   `json:"createdAt"`
	ComputedAt          string   `json:"computedAt"`
	Title               *string  `json:"title"`
	AuthorID            *string  `json:"authorID"`
	UrlSlug             *string  `json:"urlSlug"`
	Content             string   `json:"content"`
	Avatar              *string  `json:"avatar"`
	DisplayName         *string  `json:"displayName"`
	Slug                *string  `json:"slug"`
}
