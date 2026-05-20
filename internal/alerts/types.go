package alerts

import "time"

type createAlertParams struct {
	Iddevice    int32 `json:"iddevice"`
	Idtypealert int32 `json:"idtypealert"`
}

type AlertResponse struct {
	ID                   int32     `json:"id"`
	Date                 time.Time `json:"date"`
	DeviceID             int32     `json:"deviceId"`
	TypeAlertID          int32     `json:"typeAlertId"`
	TypeAlertDescription string    `json:"typeAlertDescription"`
}
