package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	CreateBlogSaga                string = "saga.create_blog"
	CreateBlog                    string = "cmd.CreateBlog"
	CreateBlogSuccess             string = "evt.CreateBlog.Success"
	CreateBlogFailed              string = "evt.CreateBlog.Failed"
	InceaseAuthorBlogCount        string = "cmd.InceaseAuthorBlogCount"
	InceaseAuthorBlogCountSuccess string = "evt.InceaseAuthorBlogCount.Success"
	InceaseAuthorBlogCountFailed  string = "evt.InceaseAuthorBlogCount.Failed"
	CreateBlogCompensation        string = "cmd.CreateBlogCompensation"
	CreateBlogCompensationSuccess string = "evt.CreateBlogCompensation.Success"
	CreateBlogCompensationFailed  string = "evt.CreateBlogCompensation.Failed"
)

var CreateBlogSagaDefinition = domain.SagaDefinition{
	Name: CreateBlogSaga,
	Steps: []domain.Step{
		{
			ActionType:     CreateBlog,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     InceaseAuthorBlogCount,
			CompensateType: CreateBlogCompensation,
			MaxRetries:     2,
		},
	},
}
