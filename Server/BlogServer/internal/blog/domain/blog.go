package domain

import (
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

type Blog struct {
	BlogID       int64           `json:"id"`
	AuthorID     string          `json:"authorID"`
	Title        string          `json:"title"`
	URLSlug      string          `json:"urlSlug"`
	ContentText  string          `json:"contentText"`
	ContentJson  json.RawMessage `json:"contentJson"`
	Status       string          `json:"status"`
	Tags         []string        `json:"tags"`
	ThumbnailUrl *string         `json:"thumbnailUrl"`
	// Images  []string
	model.Audit
}

type BlogWithAuthorData struct {
	BlogID       int64           `json:"blogID"`
	Author       AuthorData      `json:"author"`
	URLSlug      string          `json:"urlSlug"`
	Title        string          `json:"title"`
	ContentText  string          `json:"contentText"`
	ContentJson  json.RawMessage `json:"contentJson"`
	ThumbnailUrl *string         `json:"thumbnailUrl"`
	LikeCount    int64           `json:"likeCount"`
	DislikeCount int64           `json:"dislikeCount"`
	UserReaction *string         `json:"userReaction"`
	Status       string          `json:"status"`
	Tags         []string        `json:"tags"`
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
	ThumbnailUrl        *string  `json:"thumbnailUrl"`
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

type WeekViewData struct {
	WeekStart string `json:"weekStart"`
	Views     int64  `json:"views"`
}

type DateViewData struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

type BlogReport struct {
	ReportID        int64  `json:"reportID"`
	BlogID          int64  `json:"blogID"`
	UserID          string `json:"userID"`
	UserDisplayName string `json:"userDisplayName"`
	Reason          string `json:"reason"`
}

type AuthorDashboardViewMetrics struct {
	TodayViews     int64 `json:"todayViews"`
	YesterdayViews int64 `json:"yesterdayViews"`
	ThisWeekViews  int64 `json:"thisWeekViews"`
	LastWeekViews  int64 `json:"lastWeekViews"`
}

type AuthorDashboardReactionMetrics struct {
	TodayLikes        int64 `json:"todayLikes"`
	TodayDislikes     int64 `json:"todayDislikes"`
	YesterdayLikes    int64 `json:"yesterdayLikes"`
	YesterdayDislikes int64 `json:"yesterdayDislikes"`
	ThisWeekLikes     int64 `json:"thisWeekLikes"`
	ThisWeekDislikes  int64 `json:"thisWeekDislikes"`
	LastWeekLikes     int64 `json:"lastWeekLikes"`
	LastWeekDislikes  int64 `json:"lastWeekDislikes"`
}
