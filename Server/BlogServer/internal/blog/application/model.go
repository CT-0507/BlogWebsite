package application

type BlogCreatedEvent struct {
	BlogID    int64
	BlogTitle string
	AuthorID  string
}

func (e BlogCreatedEvent) EventName() string {
	return "blog.created"
}
