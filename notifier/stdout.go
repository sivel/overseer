package notifier

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sivel/overseer/status"
)

type StdoutConfig struct {
	Name string
	Type string
}

type Stdout struct {
	config *StdoutConfig
}

func NewStdout(conf []byte, filename string) Notifier {
	notifier := new(Stdout)
	var config StdoutConfig
	json.Unmarshal(conf, &config)
	notifier.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	return notifier
}

func (n *Stdout) Name() string {
	return n.config.Name
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
