package http

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	URLSlug string `json:"urlSlug" validate:"required,max=500"`
	Content string `json:"content" validate:"required"`
}

type CreateCommentRequest struct {
	ActorType       string  `json:"actor_type" validate:"required"`
	Content         string  `json:"content" validate:"required"`
	BlogID          int64   `json:"blog_id" validate:"required"`
	ParentCommentID *string `json:"parent_comment_id"`
	RootCommentID   string  `json:"root_comment_id" validate:"required"`
	Depth           int16   `json:"depth" validate:"required,min=0,max=2"`
}

type CreateBlogReactionRequest struct {
	Type string `json:"type" validate:"required"`
}
