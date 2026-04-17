package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

var DeleteAuthorSaga = domain.SagaDefinition{
	Name: "create_author_saga",
	Steps: []domain.Step{
		{
			ActionType:     "DeleteAuthor",
			CompensateType: "RestoreAuthor",
			MaxRetries:     2,
		},
		{
			ActionType:     "DeleteBlogAuthorCache",
			CompensateType: "",
			MaxRetries:     2,
		},
	},
}
