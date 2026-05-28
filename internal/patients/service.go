package patients

import (
	"context"
	"errors"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrInvalidGender = errors.New("gender must be exactly one character")
)

type Service interface {
	ListPatients(ctx context.Context) ([]repo.Patient, error)
	CreatePatient(ctx context.Context, tempPatientParams repo.CreatePatientParams) (repo.Patient, error)
	ListAlertsByPatient(ctx context.Context, patientID int64) ([]repo.ListAlertsByPatientRow, error)
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

func (s *svc) ListPatients(ctx context.Context) ([]repo.Patient, error) {
	return s.repoQuerier.ListPatients(ctx)
}

func (s *svc) CreatePatient(ctx context.Context, tempPatientParams repo.CreatePatientParams) (repo.Patient, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Patient{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repoQueries.WithTx(tx)

	if len([]rune(tempPatientParams.Gender)) != 1 {
		return repo.Patient{}, ErrInvalidGender
	}

	//create patient
	patient, err := qtx.CreatePatient(ctx, repo.CreatePatientParams(tempPatientParams))
	if err != nil {
		return repo.Patient{}, err
	}

	//commit transaction
	if err := tx.Commit(ctx); err != nil {
		return repo.Patient{}, err
	}

	return patient, nil
}

func (s *svc) ListAlertsByPatient(ctx context.Context, patientID int64) ([]repo.ListAlertsByPatientRow, error) {
	rows, err := s.repoQuerier.ListAlertsByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	resp := make([]repo.ListAlertsByPatientRow, len(rows))

	for i, r := range rows {
		resp[i] = repo.ListAlertsByPatientRow{
			ID:                   r.ID,
			Date:                 r.Date,
			Iddevice:             r.Iddevice,
			Idtypealert:          r.Idtypealert,
			AlertTypeDescription: r.AlertTypeDescription,
		}
	}

	return resp, nil
}
