package type_alerts

import (
	"context"
	"errors"
	"log"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrAlertTypeAlreadyExists = errors.New("alert type already exists")
)

type Service interface {
	ListTypeAlerts(ctx context.Context) ([]repo.TypeAlert, error)
	CreateAlertType(ctx context.Context, tempTypeAlert repo.TypeAlert) (repo.TypeAlert, error)
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

func (s *svc) CreateAlertType(ctx context.Context, tempTypeAlert repo.TypeAlert) (repo.TypeAlert, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.TypeAlert{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repoQueries.WithTx(tx)

	//check if theres identical ids
	tempID, err := qtx.GetTypeAlertByID(ctx, tempTypeAlert.ID)
	log.Println(tempID)
	if tempID.ID == tempTypeAlert.ID {
		return repo.TypeAlert{}, ErrAlertTypeAlreadyExists
	}

	//create type alert
	alert, err := qtx.CreateTypeAlert(ctx, repo.CreateTypeAlertParams(tempTypeAlert))
	if err != nil {
		return repo.TypeAlert{}, err
	}

	//commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.TypeAlert{}, err
	}

	return alert, nil
}
