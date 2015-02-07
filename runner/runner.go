package runner

import (
	"github.com/sivel/overseer/monitor"
	"github.com/sivel/overseer/notifier"
	"github.com/sivel/overseer/status"
)

type Runner struct {
	StatusChan chan *status.Status
	Monitors   []monitor.Monitor
	Notifiers  []notifier.Notifier
}

func NewRunner(monitors []monitor.Monitor, notifiers []notifier.Notifier) *Runner {
	runner := &Runner{
		StatusChan: make(chan *status.Status),
		Monitors:   monitors,
		Notifiers:  notifiers,
	}
	return runner
}

func (r *Runner) Loop() {
	for _, monitor := range r.Monitors {
		go monitor.Watch(r.StatusChan)
	}

	for {
		stat := <-r.StatusChan
		if !notifier.ShouldNotify(stat) {
			continue
		}
		for _, notifier := range r.Notifiers {
			notifier.Notify(stat)
		}
	}
}
