package alerts

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/chi/v5"
)

type mockService struct {
	listAlertsFn    func(ctx context.Context) ([]repo.ListAlertsRow, error)
	createAlertFn   func(ctx context.Context, tempAlert createAlertParams) (repo.Alert, error)
	listByPatientFn func(ctx context.Context, patientID int64) ([]AlertResponse, error)
}

func (m *mockService) ListAlerts(ctx context.Context) ([]repo.ListAlertsRow, error) {
	return m.listAlertsFn(ctx)
}

func (m *mockService) CreateAlert(ctx context.Context, tempAlert createAlertParams) (repo.Alert, error) {
	return m.createAlertFn(ctx, tempAlert)
}

func (m *mockService) ListAlertsByPatient(ctx context.Context, patientID int64) ([]AlertResponse, error) {
	return m.listByPatientFn(ctx, patientID)
}

// -------------------- ListAlerts --------------------

func TestListAlerts_Success(t *testing.T) {
	svc := &mockService{
		listAlertsFn: func(ctx context.Context) ([]repo.ListAlertsRow, error) {
			return []repo.ListAlertsRow{}, nil
		},
	}

	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()

	h.ListAlerts(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestListAlerts_Error(t *testing.T) {
	svc := &mockService{
		listAlertsFn: func(ctx context.Context) ([]repo.ListAlertsRow, error) {
			return nil, errors.New("db error")
		},
	}

	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	w := httptest.NewRecorder()

	h.ListAlerts(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Result().StatusCode)
	}
}

// -------------------- CreateAlert --------------------

func TestCreateAlert_Success(t *testing.T) {
	svc := &mockService{
		createAlertFn: func(ctx context.Context, tempAlert createAlertParams) (repo.Alert, error) {
			return repo.Alert{}, nil
		},
	}

	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.CreateAlert(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Result().StatusCode)
	}
}

func TestCreateAlert_BadRequest(t *testing.T) {
	svc := &mockService{}

	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodPost, "/alerts", bytes.NewBufferString(`invalid-json`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	h.CreateAlert(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}

// -------------------- ListAlertsByPatient --------------------

func TestListAlertsByPatient_Success(t *testing.T) {
	svc := &mockService{
		listByPatientFn: func(ctx context.Context, patientID int64) ([]AlertResponse, error) {
			if patientID != 123 {
				t.Fatalf("expected patientID 123, got %d", patientID)
			}
			return []AlertResponse{}, nil
		},
	}

	h := NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/patients/123/alerts", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("patient_id", "123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	h.ListAlertsByPatient(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestListAlertsByPatient_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	req := httptest.NewRequest(http.MethodGet, "/patients/abc/alerts", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("patient_id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	h.ListAlertsByPatient(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Result().StatusCode)
	}
}
