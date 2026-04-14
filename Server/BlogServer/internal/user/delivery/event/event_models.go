package event

type CreateNotificationsEventContext struct {
	AuthorID string
}

type CreateNotificationsEventPayload struct {
	BlogID           int64
	UrlSlug          string
	AuthorID         string
	AuthorName       string
	AuthorSlug       string
	UserID           string
	FollowerIDs      []string
	TruncatedTitle   string
	TruncatedContent string
}
