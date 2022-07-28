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

func NewMailgun(conf []byte, filename string) Notifier {
	notifier := new(Mailgun)

	var config MailgunConfig
	json.Unmarshal(conf, &config)
	notifier.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	if config.Domain == "" {
		log.Fatalf("Mailgun domain not provided: %s", filename)
	}

	if config.APIKey == "" {
		log.Fatalf("Mailgun API Key not provided: %s", filename)
	}

	if config.From == "" {
		log.Fatalf("Mailgun from address not provided: %s", filename)
	}

	if len(config.To) == 0 {
		log.Fatalf("Mailgun to address list not provided: %s", filename)
	}

	return notifier
}

func (n *Mailgun) Name() string {
	return n.config.Name
}

func (n *Mailgun) Type() string {
	return n.config.Type
}

func (n *Mailgun) Notify(stat *status.Status) {
	mg := mailgun.NewMailgun(n.config.Domain, n.config.APIKey)
	m := mg.NewMessage(
		n.config.From,
		fmt.Sprintf("[%s] %s", status.StateString(stat), stat.MonitorName),
		fmt.Sprintf(
			"%s [%dms] [%s]",
			stat.Message,
			stat.CheckDuration/1000000,
			time.Since(stat.StartOfCurrentStatus),
		),
		n.config.To...,
	)

	for i := 1; i < 60; i++ {
		_, _, err := mg.Send(m)

		if err != nil {
			log.Print("Mailgun notifier: unable to send message")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}
