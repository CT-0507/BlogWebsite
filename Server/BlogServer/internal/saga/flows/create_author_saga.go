package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	CreateAuthorSaga                string = "saga.create_author"
	CreateAuthor                    string = "cmd.CreateAuthor"
	CreateAuthorSuccess             string = "evt.CreateAuthor.Success"
	CreateAuthorFailed              string = "evt.CreateAuthor.Failed"
	CreateBlogAuthorCache           string = "cmd.CreateBlogAuthorCache"
	CreateBlogAuthorCacheSuccess    string = "evt.CreateBlogAuthorCache.Success"
	CreateBlogAuthorCacheFailed     string = "evt.CreateBlogAuthorCache.Failed"
	CreateAuthorCompensation        string = "cmd.CreateAuthorCompensation"
	CreateAuthorCompensationSuccess string = "evt.CreateAuthorCompensation.Success"
	CreateAuthorCompensationFailed  string = "evt.CreateAuthorCompensation.Failed"
)

var CreateAuthorSagaDefinition = domain.SagaDefinition{
	Name: CreateAuthorSaga,
	Steps: []domain.Step{
		{
			ActionType:     CreateAuthor,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     CreateBlogAuthorCache,
			CompensateType: CreateAuthorCompensation,
			MaxRetries:     2,
		},
	},
}
