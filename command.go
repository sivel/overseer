package main

import (
	"github.com/sivel/overseer/config"
	"github.com/sivel/overseer/runner"
)

func main() {
	monitors, notifiers := config.ParseConfig()
	run := runner.NewRunner(monitors, notifiers)
	run.Loop()
}
