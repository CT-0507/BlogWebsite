package repository

import "context"

type TagRepository interface {
	// Tags
	UpsertTags(c context.Context, blogID int64, name []string) error
}
