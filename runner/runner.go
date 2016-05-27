package runner

import (
	"log"

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
	if len(r.Monitors) == 0 {
		log.Fatalf("No monitors are configured. Exiting...")
	}

	for _, monitor := range r.Monitors {
		go monitor.Watch(r.StatusChan)
	}

	for {
		stat := <-r.StatusChan
		go func(stat *status.Status) {
			if notifier.ShouldNotify(stat) {
				for _, n := range r.Notifiers {
					go func(stat *status.Status, n notifier.Notifier) {
						if notifier.NotifierMatch(stat, n) {
							n.Notify(stat)
						}
					}(stat, n)
				}
			}
		}(stat)
	}
}
