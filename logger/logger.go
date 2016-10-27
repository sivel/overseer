package logger

import (
	"errors"
	"fmt"

	"github.com/sivel/overseer/status"
)

type Logger interface {
	Log(*status.Status)
	Name() string
	Type() string
}

type NewLogger func([]byte, string) Logger

func GetLogger(loggerType string) (NewLogger, error) {
	switch loggerType {
	case "stderr":
		return NewStderr, nil
	case "mongodb":
		return NewMongoDB, nil
	case "elasticsearch":
		return NewElasticsearch, nil
	}
	return nil, errors.New(fmt.Sprintf("Unsuppported logger type: %s", loggerType))
}

func LoggerMatch(stat *status.Status, l Logger) bool {
	if len(stat.Loggers) == 0 {
		return true
	} else {
		for _, loggerName := range stat.Loggers {
			if l.Name() == loggerName || l.Type() == loggerName {
				return true
			}
		}
	}
	return false
}
