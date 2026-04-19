package flows

import "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"

const (
	DeleteUserSaga                          string = "saga.delete_user"
	DeleteUser                              string = "cmd.DeleteUser"
	DeleteUserSuccess                       string = "evt.DeleteUser.Success"
	DeleteUserFailed                        string = "evt.DeleteUser.Failed"
	CleanUpAuthorProfile                    string = "cmd.CleanUpAuthorProfile"
	CleanUpAuthorProfileSuccess             string = "evt.CleanUpAuthorProfile.Success"
	CleanUpAuthorProfileFailed              string = "evt.CleanUpAuthorProfile.Failed"
	DeleteUserCompensation                  string = "cmd.DeleteUserCompensation"
	DeleteUserCompensationSuccess           string = "evt.DeleteUserCompensation.Success"
	DeleteUserCompensationFailed            string = "evt.DeleteUserCompensation.Failed"
	CleanUpAuthorProfileCompensation        string = "cmd.CleanUpAuthorProfileCompensation"
	CleanUpAuthorProfileCompensationSuccess string = "evt.CleanUpAuthorProfileCompensation.Success"
	CleanUpAuthorProfileCompensationFailed  string = "evt.CleanUpAuthorProfileCompensation.Failed"
)

var DeleteUserSagaDefinition = &domain.SagaDefinition{
	Name: DeleteUserSaga,
	Steps: []domain.Step{
		{
			ActionType:     DeleteUser,
			CompensateType: "",
			MaxRetries:     2,
		},
		{
			ActionType:     CleanUpAuthorProfile,
			CompensateType: DeleteUserCompensation,
			MaxRetries:     2,
		},
		{
			ActionType:     DeleteBlogAuthorCache,
			CompensateType: CleanUpAuthorProfileCompensation,
			MaxRetries:     2,
		},
	},
}
