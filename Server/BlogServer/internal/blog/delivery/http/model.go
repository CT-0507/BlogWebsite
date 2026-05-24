package http

type CreateBlogRequest struct {
	Title       string   `form:"title" validate:"required"`
	URLSlug     string   `form:"urlSlug" validate:"required,max=500"`
	ContentText string   `form:"contentText" validate:"required"`
	ContentJson string   `form:"contentJson" validate:"required"`
	Tags        []string `form:"tags" validate:"max=5"`
}

type CreateCommentRequest struct {
	ActorType       string  `json:"actorType" validate:"required"`
	Content         string  `json:"content" validate:"required"`
	ParentCommentID *string `json:"parentCommentId"`
	RootCommentID   *string `json:"rootCommentId"`
	Depth           int16   `json:"depth" validate:"min=0,max=2"`
}

type CreateBlogReactionRequest struct {
	Type string `json:"type" validate:"required"`
}

type CreateCommentReactionRequest struct {
	Type string `json:"type" validate:"required"`
}

type UpdateCommentContentRequest struct {
	Content string `json:"content" validate:"required,max=2000"`
}

type GetBlogFilter struct {
	Title      *string `form:"title" validate:"omitempty,max=50"`
	Content    *string `form:"content" validate:"omitempty,max=100"`
	AuthorName *string `form:"authorName" validate:"omitempty,max=50"`
	Page       *int32  `form:"page"  validate:"omitempty"`
	Limit      *int32  `form:"limit" validate:"omitempty"`
	SortBy     *string `form:"sortBy" validate:"omitempty"`
	SortDir    *string `form:"sortDir" validate:"omitempty"`
}

type GetBlogRankingFilter struct {
	Type    *string `form:"type" validate:"omitempty,max=10"`
	Page    *int32  `form:"page"  validate:"omitempty,min=1"`
	Limit   *int32  `form:"limit" validate:"omitempty,min=1"`
	SortBy  *string `form:"sortBy" validate:"omitempty"`
	SortDir *string `form:"sortDir" validate:"omitempty"`
}

type CreateBlogReportRequest struct {
	Reason string `json:"reason" validate:"required,min=8,max=1000"`
}

type GetBlogReportsForAuthorResponse struct {
	ReportID int64  `json:"reportID"`
	BlogID   int64  `json:"blogID"`
	Reason   string `json:"reason"`
}
