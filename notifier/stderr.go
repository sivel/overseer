package notifier

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sivel/overseer/status"
)

type StderrConfig struct {
	Name string
	Type string
}

type Stderr struct {
	config *StderrConfig
}

func NewStderr(conf []byte, filename string) Notifier {
	notifier := new(Stderr)
	var config StderrConfig
	json.Unmarshal(conf, &config)
	notifier.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	return notifier
}

func (n *Stderr) Name() string {
	return n.config.Name
}

func (n *Stderr) Notify(stat *status.Status) {
	log.Printf(
		"[%s] %s: %s [%dms] [%s]\n",
		stateString(stat),
		stat.MonitorName,
		stat.Message,
		stat.CheckDuration/1000000,
		time.Since(stat.StartOfCurrentStatus),
	)
}
