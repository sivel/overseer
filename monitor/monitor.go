package monitor

import (
	"errors"
	"fmt"
	"time"

	"github.com/sivel/overseer/status"
)

type Monitor interface {
	Check()
	Watch(chan *status.Status)
}

type NewMonitor func([]byte, string) Monitor

func GetMonitor(monitorType string) (NewMonitor, error) {
	switch monitorType {
	case "http-status":
		return NewHTTPStatus, nil
	case "http-content":
		return NewHTTPContent, nil
	case "connect":
		return NewConnect, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsuppported monitor type: %s", monitorType))
}

func checkChanged(current int, last int, startOfLastStatus time.Time) (bool, time.Time) {
	var start time.Time = startOfLastStatus
	var changed bool = false
	if current != last {
		changed = true
		start = time.Now()
	}
	return changed, start
}
