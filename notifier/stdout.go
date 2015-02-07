package notifier

import (
	"log"
	"time"

	"github.com/sivel/overseer/status"
)

type StdoutConfig struct {
	Type          string
	Someting      int
	SomethingElse int
}

type Stdout struct {
	config *StdoutConfig
}

func NewStdout(conf map[string]interface{}) Notifier {
	notifier := new(Stdout)
	notifier.config = &StdoutConfig{}
	return notifier
}

func (n *Stdout) Notify(stat *status.Status) {
	log.Printf(
		"[%s] %s: %s [%dms] [%s]\n",
		stateString(stat),
		stat.MonitorName,
		stat.Message,
		stat.CheckDuration/1000000,
		time.Since(stat.StartOfCurrentStatus),
	)
}
