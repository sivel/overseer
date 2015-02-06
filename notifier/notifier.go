package notifier

import (
	"errors"
	"fmt"

	"github.com/sivel/overseer/status"
)

type Notifier interface {
	Notify(*status.Status)
}

type NewNotifier func(map[string]interface{}) Notifier

func GetNotifier(notifierType string) (NewNotifier, error) {
	switch notifierType {
	case "stdout":
		return NewStdout, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsuppported notifier type: %s", notifierType))
}
