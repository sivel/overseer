package notifier

import (
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go"
	"github.com/sivel/overseer/status"
)

type MailgunConfig struct {
	Domain string
	APIKey string
	From   string
	To     []string
}

type Mailgun struct {
	config *MailgunConfig
}

func NewMailgun(conf map[string]interface{}) Notifier {
	notifier := new(Mailgun)

	var domain string
	if domainInterface, ok := conf["domain"]; ok {
		domain = domainInterface.(string)
	} else {
		log.Fatal("Mailgun domain not provided")
	}

	var apiKey string
	if apiKeyInterface, ok := conf["apikey"]; ok {
		apiKey = apiKeyInterface.(string)
	} else {
		log.Fatal("Mailgun API Key not provided")
	}

	var from string
	if fromInterface, ok := conf["from"]; ok {
		from = fromInterface.(string)
	} else {
		log.Fatal("Mailgun from address not provided")
	}

	var to []string
	if toInterface, ok := conf["to"]; ok {
		for _, addr := range toInterface.([]interface{}) {
			to = append(to, addr.(string))
		}
	} else {
		log.Fatal("Mailgun to address list not provided")
	}

	notifier.config = &MailgunConfig{
		Domain: domain,
		APIKey: apiKey,
		From:   from,
		To:     to,
	}
	return notifier
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
