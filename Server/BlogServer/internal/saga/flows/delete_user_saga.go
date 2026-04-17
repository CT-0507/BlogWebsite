package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

var DeleteUserSaga = domain.SagaDefinition{
	Name: "create_blog_saga",
	Steps: []domain.Step{
		{
			ActionType:     "DeleteUser",
			CompensateType: "DeleteUserCompensation",
			MaxRetries:     2,
		},
		{
			ActionType:     "DeleteAuthorProfile",
			CompensateType: "DeleteAuthorProfileCompensation",
			MaxRetries:     2,
		},
		{
			ActionType:     "DeleteBlogAuthorCache",
			CompensateType: "",
			MaxRetries:     2,
		},
	},
}
