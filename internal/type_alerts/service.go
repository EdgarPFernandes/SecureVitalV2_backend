package type_alerts

import (
	"context"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	ListTypeAlerts(ctx context.Context) ([]repo.TypeAlert, error)
}

type svc struct {
	repoQuerier repo.Querier
	repoQueries *repo.Queries
	db          *pgx.Conn
}

func NewService(repoQuerier repo.Querier, repoQueries *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repoQuerier: repoQuerier,
		repoQueries: repoQueries,
		db:          db,
	}
}

func (s *svc) ListTypeAlerts(ctx context.Context) ([]repo.TypeAlert, error) {
	return s.repoQuerier.ListTypeAlerts(ctx)
}
