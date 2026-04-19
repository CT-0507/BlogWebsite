package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	DeleteAuthorSaga                string = "saga.delete_author"
	DeleteAuthor                    string = "cmd.DeleteAuthor"
	DeleteAuthorSuccess             string = "evt.DeleteAuthor.Success"
	DeleteAuthorFailed              string = "evt.DeleteAuthor.Failed"
	DeleteBlogAuthorCache           string = "cmd.DeleteBlogAuthorCache"
	DeleteBlogAuthorCacheSuccess    string = "evt.DeleteBlogAuthorCache.Success"
	DeleteBlogAuthorCacheFailed     string = "evt.DeleteBlogAuthorCache.Failed"
	DeleteAuthorCompensation        string = "cmd.DeleteAuthorCompensation"
	DeleteAuthorCompensationSuccess string = "evt.DeleteAuthorCompensation.Success"
	DeleteAuthorCompensationFailed  string = "evt.DeleteAuthorCompensation.Failed"
)

var DeleteAuthorSagaDefinition = domain.SagaDefinition{
	Name: DeleteAuthorSaga,
	Steps: []domain.Step{
		{
			ActionType:     DeleteAuthor,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     DeleteBlogAuthorCache,
			CompensateType: DeleteAuthorCompensation,
			MaxRetries:     2,
		},
	},
}
