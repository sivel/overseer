package monitor

import (
	"errors"
	"fmt"

	"github.com/sivel/overseer/status"
)

type Monitor interface {
	Check()
	Watch(chan *status.Status)
}

type NewMonitor func(map[string]interface{}) Monitor

func GetMonitor(monitorType string) (NewMonitor, error) {
	switch monitorType {
	case "http-status":
		return NewHTTPStatus, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsuppported notifier type: %s", monitorType))
}
