package status

import "time"

const (
	DOWN    int = 0
	UP      int = 1
	UNKNOWN int = 2
)

type Status struct {
	MonitorName          string
	Current              int
	Last                 int
	NotificationInterval time.Duration
	StartOfCurrentStatus time.Time
	LastNotification     time.Time
	CheckDuration        int64
	Message              string
	Time                 time.Time
	Notifiers            []string
	Loggers              []string
}

func NewStatus(
	monitorName string, current int, last int,
	notificationInterval time.Duration, startOfCurrentStatus time.Time,
	lastNotification time.Time, checkDuration int64, message string,
	notifiers []string, loggers []string,
) *Status {
	return &Status{
		MonitorName:          monitorName,
		Current:              current,
		Last:                 last,
		NotificationInterval: notificationInterval,
		StartOfCurrentStatus: startOfCurrentStatus,
		LastNotification:     lastNotification,
		CheckDuration:        checkDuration,
		Message:              message,
		Time:                 time.Now(),
		Notifiers:            notifiers,
		Loggers:              loggers,
	}
}

func StateString(stat *Status) string {
	switch stat.Current {
	case UP:
		return "ok"
	case DOWN:
		return "critical"
	case UNKNOWN:
		return "unknown"
	}
	return "unknown"
}
