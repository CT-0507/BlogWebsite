package http

type CreateAuthorRequest struct {
	DisplayName string `json:"displayName" validate:"required"`
	Bio         string `json:"bio" validate:"max=1000"`
	Avatar      string `json:"avatar" validate:"max=150"`
	Slug        string `json:"slug" validate:"required,min=1,max=150"`
	SocialLink  string `json:"socialLink" validate:"max=200"`
	Email       string `json:"email" validate:"omitempty,email,max=200"`
}

type UpdateAuthorSlugRequest struct {
	Slug string `json:"slug" validate:"required,max=150"`
}

type UpdateAuthorStatusRequest struct {
	Status string `json:"status" validate:"required,max=15"`
}
