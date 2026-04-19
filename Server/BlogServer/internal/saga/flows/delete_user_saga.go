package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	DeleteUserSaga                string = "saga.delete_user"
	DeleteUser                    string = "cmd.DeleteUser"
	DeleteUserSuccess             string = "evt.DeleteUser.Success"
	DeleteUserFailed              string = "evt.DeleteUser.Failed"
	DeleteAuthorProfile           string = "cmd.DeleteAuthorProfile"
	DeleteAuthorProfileSuccess    string = "evt.DeleteAuthorProfile.Success"
	DeleteAuthorProfileFailed     string = "evt.DeleteAuthorProfile.Failed"
	DeleteUserCompensation        string = "cmd.DeleteUserCompensation"
	DeleteUserCompensationSuccess string = "evt.DeleteUserCompensation.Success"
	DeleteUserCompensationFailed  string = "evt.DeleteUserCompensation.Failed"
)

var DeleteUserSagaDefinition = domain.SagaDefinition{
	Name: DeleteUserSaga,
	Steps: []domain.Step{
		{
			ActionType:     DeleteUser,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     DeleteAuthorProfile,
			CompensateType: DeleteUserCompensation,
			MaxRetries:     2,
		},
		{
			ActionType:     DeleteBlogAuthorCache,
			CompensateType: DeleteAuthorCompensation,
			MaxRetries:     2,
		},
	},
}
