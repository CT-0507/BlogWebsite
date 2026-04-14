package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

var CreateBlogSaga = domain.SagaDefinition{
	Name: "create_blog_saga",
	Steps: []domain.Step{
		{
			ActionType:     "CreateBlog",
			CompensateType: "DeleteBlog",
			MaxRetries:     2,
		},
		{
			ActionType:     "InceaseAuthorBlogCount",
			CompensateType: "DecreaseAuthorBlogCount",
			MaxRetries:     2,
		},
	},
}
