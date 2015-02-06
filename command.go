package main

import (
	"fmt"

	"github.com/sivel/overseer/config"
	"github.com/sivel/overseer/runner"
)

func main() {
	fmt.Println("overseer")
	monitors, notifiers := config.ParseConfig()
	run := runner.NewRunner(monitors, notifiers)
	run.Loop()
}
