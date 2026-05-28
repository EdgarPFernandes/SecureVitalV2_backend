package devices

import (
	"context"
	"errors"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrPatientNotFound = errors.New("patient not found")
)

type Service interface {
	ListDevices(ctx context.Context) ([]repo.ListDevicesRow, error)
	CreateDevice(ctx context.Context, tempDevice repo.CreateDeviceParams) (repo.Device, error)
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

func (s *svc) ListDevices(ctx context.Context) ([]repo.ListDevicesRow, error) {
	return s.repoQuerier.ListDevices(ctx)
}

func (s *svc) CreateDevice(ctx context.Context, tempDevice repo.CreateDeviceParams) (repo.Device, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Device{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repoQueries.WithTx(tx)

	//look for patient id
	_, err = qtx.CreateDevice(ctx, tempDevice)
	if err != nil {
		return repo.Device{}, ErrPatientNotFound
	}

	//create device
	device, err := qtx.CreateDevice(ctx, tempDevice)
	if err != nil {
		return repo.Device{}, err
	}

	//commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Device{}, err
	}

	return device, nil
}
