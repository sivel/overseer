package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go"
	"github.com/sivel/overseer/status"
)

type MailgunConfig struct {
	Name   string
	Type   string
	Domain string
	APIKey string
	From   string
	To     []string
}

type Mailgun struct {
	config *MailgunConfig
}

func NewMailgun(conf []byte) Notifier {
	notifier := new(Mailgun)

	var config MailgunConfig
	err := json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	} else {
		notifier.config = &config
	}

	if config.Name == "" {
		config.Name = config.Type
	}

	if config.Domain == "" {
		log.Fatal("Mailgun domain not provided")
	}

	if config.APIKey == "" {
		log.Fatal("Mailgun API Key not provided")
	}

	if config.From == "" {
		log.Fatal("Mailgun from address not provided")
	}

	if len(config.To) == 0 {
		log.Fatal("Mailgun to address list not provided")
	}

	return notifier
}

func (n *Mailgun) Name() string {
	return n.config.Name
}

func (n *Mailgun) Notify(stat *status.Status) {
	mg := mailgun.NewMailgun(n.config.Domain, n.config.APIKey, "")
	m := mg.NewMessage(
		n.config.From,
		fmt.Sprintf("[%s] %s", stateString(stat), stat.MonitorName),
		fmt.Sprintf(
			"%s [%dms] [%s]",
			stat.Message,
			stat.CheckDuration/1000000,
			time.Since(stat.StartOfCurrentStatus),
		),
		n.config.To...,
	)

	_, _, err := mg.Send(m)

	if err != nil {
		log.Print("Mailgun notifier: unable to connect to send message")
	}
}
