package runner

import (
	"log"

	"github.com/sivel/overseer/logger"
	"github.com/sivel/overseer/monitor"
	"github.com/sivel/overseer/notifier"
	"github.com/sivel/overseer/status"
)

type Runner struct {
	StatusChan chan *status.Status
	Monitors   []monitor.Monitor
	Notifiers  []notifier.Notifier
	Loggers    []logger.Logger
}

func NewRunner(monitors []monitor.Monitor, notifiers []notifier.Notifier, loggers []logger.Logger) *Runner {
	runner := &Runner{
		StatusChan: make(chan *status.Status),
		Monitors:   monitors,
		Notifiers:  notifiers,
		Loggers:    loggers,
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
			go func(stat *status.Status) {
				for _, l := range r.Loggers {
					go func(stat *status.Status, l logger.Logger) {
						if logger.LoggerMatch(stat, l) {
							l.Log(stat)
						}
					}(stat, l)
				}
			}(stat)

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
