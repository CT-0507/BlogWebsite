package application

type BlogCreatedEvent struct {
	BlogID    int64
	BlogTitle string
}

func (e BlogCreatedEvent) EventName() string {
	return "blog.created"
}
