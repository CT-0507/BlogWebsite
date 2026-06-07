package infrastructure

import (
	"encoding/json"

	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/blog/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/model"
)

type DbBlog struct {
	BlogID       int64           `json:"blog_id"`
	AuthorID     string          `json:"author_id"`
	Title        string          `json:"title"`
	URLSlug      string          `json:"url_slug"`
	ContentText  string          `json:"content_text"`
	ContentJson  json.RawMessage `json:"content_json"`
	Status       string          `json:"status"`
	ThumbnailUrl *string         `json:"thumbnailUrl"`
	LikeCount    int64           `json:"likeCount"`
	DislikeCount int64           `json:"dislikeCount"`
	UserReaction *string         `json:"userReaction"`
	// Images  []string
	model.Audit
}

func (m *DbBlog) toDomainBlog() *domain.Blog {
	return &domain.Blog{
		BlogID:       m.BlogID,
		Title:        m.Title,
		URLSlug:      m.URLSlug,
		ContentJson:  m.ContentJson,
		ContentText:  m.ContentText,
		ThumbnailUrl: m.ThumbnailUrl,
		Status:       m.Status,
	}
}
