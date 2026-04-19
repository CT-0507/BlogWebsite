package http

type CreateAuthorRequest struct {
	DisplayName string `form:"displayName" validate:"required"`
	Bio         string `form:"bio" validate:"max=1000"`
	Avatar      string `form:"avatar" validate:"max=150"`
	Slug        string `form:"slug" validate:"required,min=1,max=150"`
	SocialLink  string `form:"socialLink" validate:"max=200"`
	Email       string `form:"email" validate:"omitempty,email,max=200"`
}

type UpdateAuthorSlugRequest struct {
	Slug string `json:"slug" validate:"required,max=150"`
}

type UpdateAuthorStatusRequest struct {
	Status string `json:"status" validate:"required,max=15"`
}
