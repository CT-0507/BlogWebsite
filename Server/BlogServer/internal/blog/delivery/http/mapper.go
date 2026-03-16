package http

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required"`
	URLSlug string `json:"urlSlug" validate:"required,max=500"`
	Content string `json:"content" validate:"required"`
}
