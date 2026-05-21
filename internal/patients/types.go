package patients

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type createPatientParams struct {
	Name             string      `json:"name"`
	Birthdate        pgtype.Date `json:"birth_date"`
	Gender           string      `json:"gender"`
	Address          pgtype.Text `json:"address"`
	EmergencyContact pgtype.Text `json:"emergency_contact"`
}
