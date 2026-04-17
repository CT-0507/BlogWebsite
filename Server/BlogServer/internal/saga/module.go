package saga

import (
	outboxrepo "github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/contracts/outboxRepo"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/domain"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/flows"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/saga/infrastructure"
	"github.com/CT-0507/BlogWebsite/Server/BlogServer/internal/shared/database"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SagaModule struct {
	Orchestrator domain.Orchestrator
}

func NewSagaModule(pool *pgxpool.Pool, txManageer database.TxManager, outboxRepo outboxrepo.OutboxRepository) *SagaModule {

	repo := infrastructure.NewSagaRepository(pool)

	// Register in-memory definition
	registry := infrastructure.NewRegistry()
	registry.Register(flows.CreateBlogSaga.Name, flows.CreateBlogSaga.Steps)
	registry.Register(flows.CreateAuthorSaga.Name, flows.CreateAuthorSaga.Steps)
	registry.Register(flows.DeleteAuthorSaga.Name, flows.DeleteAuthorSaga.Steps)
	registry.Register(flows.DeleteUserSaga.Name, flows.DeleteUserSaga.Steps)

	orchestrator := infrastructure.NewOrchestrator(registry, txManageer, repo, outboxRepo)
	return &SagaModule{
		Orchestrator: orchestrator,
	}
}
