package status

import "time"

const (
	DOWN    int = 0
	UP      int = 1
	UNKNOWN int = 2
)

type Status struct {
	Current              int
	Last                 int
	StartOfCurrentStatus time.Time
	LastNotification     time.Time
	CheckDuration        int64
	Message              string
	Time                 time.Time
}

func NewStatus(current int, last int, startOfCurrentStatus time.Time, lastNotification time.Time, checkDuration int64, message string) *Status {
	return &Status{
		Current:              current,
		Last:                 last,
		StartOfCurrentStatus: startOfCurrentStatus,
		LastNotification:     lastNotification,
		CheckDuration:        checkDuration,
		Message:              message,
		Time:                 time.Now(),
	}
}
