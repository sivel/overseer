package main

import (
	"runtime"

	"github.com/sivel/overseer/config"
	"github.com/sivel/overseer/runner"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	monitors, notifiers, loggers := config.ParseConfig()
	run := runner.NewRunner(monitors, notifiers, loggers)
	run.Loop()
}
