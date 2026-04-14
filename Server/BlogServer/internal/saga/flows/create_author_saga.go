package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

var CreateAuthorSaga = domain.SagaDefinition{
	Name: "create_author_saga",
	Steps: []domain.Step{
		{
			ActionType:     "CreateAuthor",
			CompensateType: "DeleteAuthor",
			MaxRetries:     2,
		},
		{
			ActionType:     "CreateBlogAuthorCache",
			CompensateType: "DeleteBlogAuthorCache",
			MaxRetries:     2,
		},
	},
}
