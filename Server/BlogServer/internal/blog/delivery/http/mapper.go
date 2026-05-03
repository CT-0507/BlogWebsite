package http

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	URLSlug string `json:"urlSlug" validate:"required,max=500"`
	Content string `json:"content" validate:"required"`
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
