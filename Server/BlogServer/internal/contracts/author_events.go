package contracts

// Author module
type AuthorCreatedEventPayload struct {
	AuthorID    string
	UserID      string
	Slug        string
	Status      string
	DisplayName string
	SocialLink  *string
	Email       *string
	Bio         *string
	Avatar      *string
	CreatedBy   *string
}

type AuthorCreatedEventContext struct {
	UserID string
	Avatar *string
}

type CreateBlogAuthorCachePayload struct {
	UserID      string
	AuthorID    string
	Slug        string
	DisplayName string
}

type CreateBlogAuthorCacheContext struct {
	AuthorID string
	Avatar   *string
}

type AuthorDeleteContext struct {
	AuthorID string
	UserID   string
	Status   string
}

type AuthorDeletePayload struct {
	AuthorID string
	UserID   string
	Status   string
	Avatar   *string
}

type AuthorDeletedEvent struct {
	AuthorID string
}

func (e AuthorDeletedEvent) EventName() string {
	return "authorIdentity.deleted"
}

type AuthorHardDeletedEvent struct {
	AuthorID string
}

func (e AuthorHardDeletedEvent) EventName() string {
	return "authorIdentity.hardDeleted"
}

type AuthorFollowedEvent struct {
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e AuthorFollowedEvent) EventName() string {
	return "authorFollower.created"
}

type AuthorUnfollowedEvent struct {
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e AuthorUnfollowedEvent) EventName() string {
	return "authorFollower.deleted"
}

type FollowCountChangedEvent struct {
	Slug        string
	AuthorID    string
	UserID      string
	IsIncrement bool
}

func (e FollowCountChangedEvent) EventName() string {
	eventType := "Increased"
	if !e.IsIncrement {
		eventType = "Decreased"
	}
	return "authorFollower.followerCount" + eventType
}

type BlogCountChangedEvent struct {
	BlogID           int64
	AuthorID         string
	UserID           string
	TruncatedTitle   string
	TruncatedContent string
	UrlSlug          string
}

func (e BlogCountChangedEvent) EventName() string {
	return "authorIdentity.blogCountIncreased"
}
