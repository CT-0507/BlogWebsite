package contracts

// Author module
type AuthorCreatedEvent struct {
	AuthorID    string
	UserID      string
	Slug        string
	DisplayName string
}

func (e AuthorCreatedEvent) EventName() string {
	return "authorIdentity.created"
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
	BlogID   int64
	AuthorID string
}

func (e BlogCountChangedEvent) EventName() string {
	return "authorIdentity.blogCountIncreased"
}
