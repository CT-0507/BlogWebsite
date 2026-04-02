package create_blog

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

var CreateBlogSaga = domain.SagaDefinition{
	Name: "create_blog_saga",
	Steps: []domain.Step{
		{
			ActionType:     "CreateBlog",
			CompensateType: "DeleteBlog",
		},
		{
			ActionType:     "InceaseAuthorBlogCount",
			CompensateType: "DecreaseAuthorBlogCount",
		},
		// Move on to normal event
		// {
		// 	ActionType:     "CreateNotifications",
		// 	CompensateType: "DeleteNotifications",
		// },

	},
}
