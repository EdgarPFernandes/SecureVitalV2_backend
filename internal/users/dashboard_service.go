package users

import (
	"context"
	"time"

	repo "github.com/EdgarPFernandes/SecureVitalV2_backend/internal/adapters/postgresql/sqlc"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userID int64) ([]PatientDashboard, error)
	GetAdminDashboard(ctx context.Context) (*AdminDashboard, error)
	GetPatientDashboard(ctx context.Context,patientID int64) (*PatientDetailsDashboard, error)
}

type PatientDashboard struct {
	ID          int64             `json:"id"`
	Name        string            `json:"name"`
	Gender      string            `json:"gender"`
	Alerts24h   int64             `json:"alerts24h"`
	TotalAlerts int64             `json:"totalAlerts"`
	Alerts7d    int64             `json:"alerts7d"`
	AlertTypes  map[string]int64  `json:"alertTypes"`
	Monthly     []int64           `json:"monthly"`
	Yearly      []int64           `json:"yearly"`
	LastAlerts  []AlertItem       `json:"lastAlerts"`
}


type AdminDashboard struct {
	TotalElderly            int64             `json:"totalElderly"`
	Alerts24h               int64             `json:"alerts24h"`

	LastWeekAlerts          []int64          `json:"lastWeekAlerts"`
	LastYearAlerts          []int64          `json:"lastYearAlerts"`

	LastWeekAlertPercentage float64           `json:"lastWeekAlertPercentage"`

	AlertTypes              map[string]int64 `json:"alertTypes"`

	TopElderly24h           []TopElderlyItem `json:"topElderly24h"`
	TopElderly7d            []TopElderlyItem `json:"topElderly7d"`
	TopElderlyTotal         []TopElderlyItem `json:"topElderlyTotal"`
}

type PatientDetailsDashboard struct {
	ID               int64               `json:"id"`
	Name             string              `json:"name"`
	Gender           string              `json:"gender"`
	BirthDate        string              `json:"birthDate"`
	Address          string              `json:"address"`
	EmergencyContact string              `json:"emergencyContact"`
	CreatedAt        string              `json:"createdAt"`

	Caregivers      []CaregiverItem     `json:"caregivers"`

	Alerts24h      int64                `json:"alerts24h"`
	Alerts7d       int64                `json:"alerts7d"`
	TotalAlerts    int64                `json:"totalAlerts"`

	Monthly         []int64             `json:"monthly"`
	Yearly          []int64             `json:"yearly"`
	AlertTypes     map[string]int64    `json:"alertTypes"`
	LastAlerts      []AlertItem 		`json:"lastAlerts"`
}

type TopElderlyItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Alerts24h   int64  `json:"alerts24h"`
	Alerts7d    int64  `json:"alerts7d"`
	TotalAlerts int64  `json:"totalAlerts"`
}


type dashboardSvc struct {
	repo repo.Querier
}

type AlertItem struct {
	ID   int32  `json:"id"`
	Date string `json:"date"`
	Type string `json:"type"`
}

type CaregiverItem struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

func NewDashboardService(repo repo.Querier) DashboardService {
	return &dashboardSvc{
		repo: repo,
	}
}

func (s *dashboardSvc) GetDashboard(ctx context.Context, userID int64) ([]PatientDashboard, error) {

	patients, err := s.repo.ListPatientsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []PatientDashboard

	for _, p := range patients {

		total, _ := s.repo.CountAlertsByPatient(ctx, p.ID)
		h24, _ := s.repo.CountAlerts24h(ctx, p.ID)
		h7d, _ := s.repo.CountAlerts7d(ctx, p.ID)

		monthlyRows, _ := s.repo.CountAlertsByMonth(ctx, p.ID)
		yearlyRows, _ := s.repo.CountAlertsByYear(ctx, p.ID)

		typesRows, _ := s.repo.ListAlertTypesByPatient(ctx, p.ID)

		// MONTHLY
		monthly := make([]int64, 12)
		for _, m := range monthlyRows {
			month := int64(m.Month)
			total := m.Total

			if month >= 1 && month <= 12 {
				monthly[month-1] = total
			}
		}

		// YEARLY
		yearly := make([]int64, 5)
		currentYear := time.Now().Year()

		for _, y := range yearlyRows {
			year := int(y.Year)
			total := y.Total

			index := currentYear - year

			if index >= 0 && index < 5 {
				yearly[4-index] = total
			}
		}

		// ALERT TYPES
		alertTypes := make(map[string]int64)
		for _, t := range typesRows {
			alertTypes[t.Description] = t.Count
		}

		// LAST ALERTS
		lastAlertsRows, _ := s.repo.ListLastAlertsByPatient(ctx, p.ID)

		lastAlerts := make([]AlertItem, 0, len(lastAlertsRows))

		for _, a := range lastAlertsRows {
			lastAlerts = append(lastAlerts, AlertItem{
				ID:   a.ID,
				Date: a.Date.Time.Format("2006-01-02 15:04"),
				Type: a.Type,
			})
		}

		result = append(result, PatientDashboard{
			ID:          p.ID,
			Name:        p.Name,
			Gender:      string(p.Gender),
			Alerts24h:   h24,
			TotalAlerts: total,
			Alerts7d:    h7d,
			AlertTypes:  alertTypes,
			Monthly:     monthly,
			Yearly:      yearly,
			LastAlerts:  lastAlerts,
		})
	}

	return result, nil
}

func (s *dashboardSvc) GetAdminDashboard(ctx context.Context) (*AdminDashboard, error) {

	// -------------------------
	// TOTALS BÁSICOS
	// -------------------------
	totalElderly, err := s.repo.CountPatients(ctx)
	if err != nil {
		return nil, err
	}

	alerts24h, err := s.repo.CountAlerts24hGlobal(ctx)
	if err != nil {
		return nil, err
	}

	// -------------------------
	// LAST 7 DAYS (ATUAL)
	// -------------------------
	weekRows, err := s.repo.CountAlertsLast7Days(ctx)
	if err != nil {
		return nil, err
	}

	dayMap := make(map[string]int64)
	for _, r := range weekRows {
		day := r.Day.Time.Format("2006-01-02")
		dayMap[day] = r.Total
	}

	lastWeek := make([]int64, 7)
	today := time.Now()

	var currentWeekTotal int64

	for i := 6; i >= 0; i-- {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		val := dayMap[day]
		lastWeek[6-i] = val
		currentWeekTotal += val
	}

	// -------------------------
	// LAST 14–7 DAYS (SEMANA ANTERIOR)
	// -------------------------
	prevRows, err := s.repo.CountAlerts7To14Days(ctx)
	if err != nil {
		return nil, err
	}

	prevMap := make(map[string]int64)
	for _, r := range prevRows {
		day := r.Day.Time.Format("2006-01-02")
		prevMap[day] = r.Total
	}

	var prevWeekTotal int64
	for i := 14; i >= 7; i-- {
		day := today.AddDate(0, 0, -i).Format("2006-01-02")
		prevWeekTotal += prevMap[day]
	}

	var lastWeekAlertPercentage float64
	if prevWeekTotal > 0 {
		lastWeekAlertPercentage =
			float64(currentWeekTotal-prevWeekTotal) / float64(prevWeekTotal) * 100
	}

	// -------------------------
	// ALERT TYPES
	// -------------------------
	typeRows, err := s.repo.CountAlertsByTypeGlobal(ctx)
	if err != nil {
		return nil, err
	}

	alertTypes := make(map[string]int64)
	for _, t := range typeRows {
		alertTypes[t.Description] = t.Total
	}

	// -------------------------
	// TOP ELDERLY
	// -------------------------
	top24Rows, err := s.repo.TopPatientsByAlerts24h(ctx)
	if err != nil {
		return nil, err
	}

	top7Rows, err := s.repo.TopPatientsByAlerts7d(ctx)
	if err != nil {
		return nil, err
	}

	topTotalRows, err := s.repo.TopPatientsByAlertsTotal(ctx)
	if err != nil {
		return nil, err
	}

	mapTop := func(rows any) []TopElderlyItem {

		switch r := rows.(type) {

		case []repo.TopPatientsByAlerts24hRow:
			res := make([]TopElderlyItem, 0, len(r))
			for _, x := range r {
				res = append(res, TopElderlyItem{
					ID:        x.ID,
					Name:      x.Name,
					Alerts24h: x.Alerts24h,
				})
			}
			return res

		case []repo.TopPatientsByAlerts7dRow:
			res := make([]TopElderlyItem, 0, len(r))
			for _, x := range r {
				res = append(res, TopElderlyItem{
					ID:       x.ID,
					Name:     x.Name,
					Alerts7d: x.Alerts7d,
				})
			}
			return res

		case []repo.TopPatientsByAlertsTotalRow:
			res := make([]TopElderlyItem, 0, len(r))
			for _, x := range r {
				res = append(res, TopElderlyItem{
					ID:          x.ID,
					Name:        x.Name,
					TotalAlerts: x.Totalalerts,
				})
			}
			return res
		}

		return nil
	}

	// -------------------------
	// LAST YEAR ALERTS
	// -------------------------
	yearRows, err := s.repo.CountAlertsLast12MonthsGlobal(ctx)
	if err != nil {
		return nil, err
	}

	lastYear := make([]int64, 12)

	for _, r := range yearRows {
		month := int(r.Month)
		total := r.Total

		if month >= 1 && month <= 12 {
			lastYear[month-1] = total
		}
	}

	// -------------------------
	// RESPONSE FINAL
	// -------------------------
	return &AdminDashboard{
		TotalElderly: totalElderly,
		Alerts24h:    alerts24h,

		LastWeekAlerts: lastWeek,
		LastYearAlerts: lastYear,

		LastWeekAlertPercentage: lastWeekAlertPercentage,

		AlertTypes: alertTypes,

		TopElderly24h:   mapTop(top24Rows),
		TopElderly7d:    mapTop(top7Rows),
		TopElderlyTotal: mapTop(topTotalRows),
	}, nil
}

func (s *dashboardSvc) GetPatientDashboard(
	ctx context.Context,
	patientID int64,
) (*PatientDetailsDashboard, error) {

	patient, err := s.repo.GetPatientDashboardInfo(
		ctx,
		patientID,
	)

	if err != nil {
		return nil, err
	}

	caregiverRows, err := s.repo.ListCaregiversByPatient(ctx, patientID)
	if err != nil {
		return nil, err
	}

	caregivers := make([]CaregiverItem, 0)

	for _, c := range caregiverRows {
		caregivers = append(caregivers, CaregiverItem{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	total, _ := s.repo.CountAlertsByPatient(ctx, patientID)
	h24, _ := s.repo.CountAlerts24h(ctx, patientID)
	h7d, _ := s.repo.CountAlerts7d(ctx, patientID)

	monthlyRows, _ := s.repo.CountAlertsByMonth(ctx, patientID)
	yearlyRows, _ := s.repo.CountAlertsByYear(ctx, patientID)

	typeRows, _ := s.repo.ListAlertTypesByPatient(ctx, patientID)

	monthly := make([]int64, 12)

	for _, m := range monthlyRows {
		if m.Month >= 1 && m.Month <= 12 {
			monthly[m.Month-1] = m.Total
		}
	}

	yearly := make([]int64, 5)

	currentYear := time.Now().Year()

	for _, y := range yearlyRows {

		index := currentYear - int(y.Year)

		if index >= 0 && index < 5 {
			yearly[4-index] = y.Total
		}
	}

	alertTypes := make(map[string]int64)

	for _, t := range typeRows {
		alertTypes[t.Description] = t.Count
	}

	lastAlertsRows, _ := s.repo.ListLastAlertsByPatient(
		ctx,
		patientID,
	)

	lastAlerts := make([]AlertItem, 0)

	for _, a := range lastAlertsRows {

		lastAlerts = append(
			lastAlerts,
			AlertItem{
				ID: a.ID,
				Date: a.Date.Time.Format(
					"2006-01-02 15:04",
				),
				Type: a.Type,
			},
		)
	}

	var address string
	if patient.Address.Valid {
		address = patient.Address.String
	}

	var emergency string
	if patient.EmergencyContact.Valid {
		emergency = patient.EmergencyContact.String
	}

	var birthDate string
	if patient.BirthDate.Valid {
		birthDate = patient.BirthDate.Time.Format("2006-01-02")
	}

	var createdAt string
	if patient.CreatedAt.Valid {
		createdAt = patient.CreatedAt.Time.Format("2006-01-02 15:04")
	}

	return &PatientDetailsDashboard{
		ID:               patient.ID,
		Name:             patient.Name,
		Gender:           string(patient.Gender),
		BirthDate:        birthDate,
		Address:          address,
		EmergencyContact: emergency,
		CreatedAt:        createdAt,

		Caregivers: caregivers,

		Alerts24h:   h24,
		Alerts7d:    h7d,
		TotalAlerts: total,
		Monthly:     monthly,
		Yearly:      yearly,
		AlertTypes:  alertTypes,
		LastAlerts:  lastAlerts,
	}, nil
}