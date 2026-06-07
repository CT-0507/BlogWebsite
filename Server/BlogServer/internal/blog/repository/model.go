package repository

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"

type ReactionTransition int

const (
	AddLike ReactionTransition = iota
	AddDislike
	LikeToDislike
	DislikeToLike
)

type BeforeAndAfterBlogUpdated struct {
	Before domain.Blog
	After  domain.Blog
}
