package logger

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

func NewStderr(conf []byte, filename string) Logger {
	logger := new(Stderr)
	var config StderrConfig
	json.Unmarshal(conf, &config)
	logger.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	return logger
}

func (n *Stderr) Name() string {
	return n.config.Name
}

func (n *Stderr) Type() string {
	return n.config.Type
}

func (n *Stderr) Log(stat *status.Status) {
	log.Printf(
		"[%s] %s: %s [%dms] [%s]\n",
		status.StateString(stat),
		stat.MonitorName,
		stat.Message,
		stat.CheckDuration/1000000,
		time.Since(stat.StartOfCurrentStatus),
	)
}
