package contracts

import "github.com/google/uuid"

type DeleteUserSagaContext struct {
	UserID    uuid.UUID
	UpdatedBy uuid.UUID
}

type DeleteUserSagaPayload struct {
	UserID    uuid.UUID
	UpdatedBy uuid.UUID
	Status    string
}

type DeleteUserContext struct {
	UserID         uuid.UUID
	PreviousStatus string
}

type DeleteUserPayload struct {
	UserID    uuid.UUID
	UpdatedBy uuid.UUID
}

type DeleteAuthorProfileContext struct {
}

type DeleteAuthorProfilePayload struct {
}
