package notifier

import (
	"errors"
	"fmt"
	"time"

	"github.com/sivel/overseer/status"
)

type Notifier interface {
	Notify(*status.Status)
}

type NewNotifier func([]byte) Notifier

func GetNotifier(notifierType string) (NewNotifier, error) {
	switch notifierType {
	case "stdout":
		return NewStdout, nil
	case "mailgun":
		return NewMailgun, nil
	case "slack":
		return NewSlack, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsuppported notifier type: %s", notifierType))
}

func ShouldNotify(stat *status.Status) bool {
	var notify bool = false
	if stat.Current == stat.Last && stat.Current == status.UP {
		notify = false
	} else if stat.Current == status.UP && stat.Last == status.UNKNOWN {
		notify = false
	} else if stat.Current == stat.Last && time.Since(stat.LastNotification) > stat.NotificationInterval {
		notify = true
	} else if stat.Current != stat.Last {
		notify = true
	}

	if notify {
		stat.LastNotification = time.Now()
	}
	return notify
}

func stateString(stat *status.Status) string {
	switch stat.Current {
	case status.UP:
		return "ok"
	case status.DOWN:
		return "critical"
	case status.UNKNOWN:
		return "unknown"
	}
	return "unknown"
}
