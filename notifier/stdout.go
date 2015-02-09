package notifier

import (
	"encoding/json"
	"log"
	"time"

	"github.com/sivel/overseer/status"
)

type StdoutConfig struct {
	Type string
}

type Stdout struct {
	config *StdoutConfig
}

func NewStdout(conf []byte) Notifier {
	notifier := new(Stdout)
	var config StdoutConfig
	err := json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	} else {
		notifier.config = &config
	}
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
