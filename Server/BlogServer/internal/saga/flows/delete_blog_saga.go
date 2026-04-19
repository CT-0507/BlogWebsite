package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	DeleteBlogSaga                 string = "saga.delete_blog"
	DeleteBlog                     string = "cmd.DeleteBlog"
	DeleteBlogSuccess              string = "evt.DeleteBlog.Success"
	DeleteBlogFailed               string = "evt.DeleteBlog.Failed"
	DecreaseAuthorBlogCount        string = "cmd.DecreaseAuthorBlogCount"
	DecreaseAuthorBlogCountSuccess string = "evt.DecreaseAuthorBlogCount.Success"
	DecreaseAuthorBlogCountFailed  string = "evt.DecreaseAuthorBlogCount.Failed"
	DeleteBlogCompensation         string = "cmd.DeleteBlogCompensation"
	DeleteBlogCompensationSuccess  string = "evt.DeleteBlogCompensation.Success"
	DeleteBlogCompensationFailed   string = "evt.DeleteBlogCompensation.Failed"
)

var DeleteBlogSagaDefinition = &domain.SagaDefinition{
	Name: DeleteBlogSaga,
	Steps: []domain.Step{
		{
			ActionType:     DeleteBlog,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     DecreaseAuthorBlogCount,
			CompensateType: DeleteBlogCompensation,
			MaxRetries:     2,
		},
	},
}
