package contracts

import "github.com/google/uuid"

type DeleteUserSagaContext struct {
	UserID    uuid.UUID
	UpdatedBy uuid.UUID
}

type DeleteUserSagaPayload struct {
	UserID    uuid.UUID
	UpdatedBy uuid.UUID
}

type DeleteUserContext struct {
}

type DeleteUserPayload struct {
}

type DeleteAuthorProfileContext struct {
}

type DeleteAuthorProfilePayload struct {
}
