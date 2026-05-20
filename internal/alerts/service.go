package alerts

import (
	"context"
	"errors"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrDeviceNotFound    = errors.New("device not found")
	ErrAlertTypeNotFound = errors.New("alert type not found")
)

type Service interface {
	ListAlerts(ctx context.Context) ([]repo.ListAlertsRow, error)
	CreateAlert(ctx context.Context, tempAlert createAlertParams) (repo.Alert, error)
	ListAlertsByPatient(ctx context.Context, patientID int64) ([]AlertResponse, error)
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

func (s *svc) ListAlerts(ctx context.Context) ([]repo.ListAlertsRow, error) {
	return s.repoQuerier.ListAlerts(ctx)
}

func (s *svc) CreateAlert(ctx context.Context, tempAlert createAlertParams) (repo.Alert, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Alert{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repoQueries.WithTx(tx)

	//look for iddevice and idtyalert
	_, err = qtx.GetDeviceByID(ctx, int64(tempAlert.Iddevice))
	if err != nil {
		return repo.Alert{}, ErrDeviceNotFound
	}

	_, err = qtx.GetTypeAlertByID(ctx, tempAlert.Idtypealert)
	if err != nil {
		return repo.Alert{}, ErrAlertTypeNotFound
	}

	//create alert
	alert, err := qtx.CreateAlert(ctx, repo.CreateAlertParams(tempAlert))
	if err != nil {
		return repo.Alert{}, err
	}

	//commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Alert{}, err
	}

	return alert, nil
}

func (s *svc) ListAlertsByPatient(ctx context.Context, patientID int64) ([]AlertResponse, error) {
	rows, err := s.repoQuerier.ListAlertsByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	resp := make([]AlertResponse, len(rows))

	for i, r := range rows {
		resp[i] = AlertResponse{
			ID:                   r.ID,
			Date:                 r.Date.Time,
			DeviceID:             r.Iddevice,
			TypeAlertID:          r.Idtypealert,
			TypeAlertDescription: r.AlertTypeDescription,
		}
	}

	return resp, nil
}
