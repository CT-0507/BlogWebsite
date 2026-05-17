package application

import (
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
)

type RequestTrackingUsecases struct {
	repo domain.SagaRepository
}

func NewRequestTrackingUsecases(
	repo domain.SagaRepository,
) *RequestTrackingUsecases {
	return &RequestTrackingUsecases{
		repo: repo,
	}
}
