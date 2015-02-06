package notifier

import (
	"fmt"

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

func (n *Stdout) Notify(status *status.Status) {
	fmt.Println("Stdout notified, status obj follows")
	fmt.Println(status)
}
