package dashboard

import (
	"context"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
)

type Service interface {
	GetAdminDashboard(ctx context.Context) (AdminDashboard, error)
}

type AdminDashboard struct {
	TotalElderly int64 `json:"totalElderly"`
	Alerts24h    int64 `json:"alerts24h"`
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) GetAdminDashboard(ctx context.Context) (AdminDashboard, error) {

	totalPatients, err := s.repo.CountPatients(ctx)
	if err != nil {
		return AdminDashboard{}, err
	}

	alerts24h, err := s.repo.CountAlerts24hGlobal(ctx)
	if err != nil {
		return AdminDashboard{}, err
	}

	return AdminDashboard{
		TotalElderly: totalPatients,
		Alerts24h:    alerts24h,
	}, nil
}