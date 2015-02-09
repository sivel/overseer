package notifier

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
	"github.com/sivel/overseer/status"
)

type SlackConfig struct {
	Name      string
	Type      string
	Token     string
	Channel   string
	ChannelID string `json:"channel_id"`
	Username  string
}

type Slack struct {
	config *SlackConfig
}

func NewSlack(conf []byte) Notifier {
	notifier := new(Slack)

	var config SlackConfig
	err := json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	} else {
		notifier.config = &config
	}

	if config.Name == "" {
		config.Name = config.Type
	}

	if config.Token == "" {
		log.Fatal("Slack token not provided")
	}
	api := slack.New(config.Token)

	if config.ChannelID == "" {
		if config.Channel == "" {
			config.Channel = "overseer"
		}
		channels, err := api.GetChannels(true)
		if err != nil {
			log.Printf("Cannot send message: %s", err)
		}
		for _, channel := range channels {
			if channel.Name == config.Channel {
				config.ChannelID = channel.Id
				break
			}
		}
		if config.ChannelID == "" {
			log.Printf("Could not locate slack channel: %s", config.Channel)
		}

		_, err = api.GetChannelInfo(config.ChannelID)
		if err != nil {
			log.Fatalf("Slack channel does not exist")
		}
	}

	if config.Username == "" {
		config.Username = "overseer"
	}

	return notifier
}

func (n *Slack) Name() string {
	return n.config.Name
}

func slackColor(stat *status.Status) string {
	switch stat.Current {
	case status.UP:
		return "good"
	case status.DOWN:
		return "danger"
	}
	return "warning"
}

func (n *Slack) Notify(stat *status.Status) {
	api := slack.New(n.config.Token)

	params := slack.PostMessageParameters{
		Username: n.config.Username,
	}
	attachment := slack.Attachment{
		Color: slackColor(stat),
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: stateString(stat),
				Value: fmt.Sprintf(
					"%s [%dms] [%s]",
					stat.Message,
					stat.CheckDuration/1000000,
					time.Since(stat.StartOfCurrentStatus),
				),
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	message := fmt.Sprintf("[%s] %s", stateString(stat), stat.MonitorName)
	_, _, err := api.PostMessage(n.config.ChannelID, message, params)

	if err != nil {
		log.Print("Slack notifier: unable to connect to send message")
	}
}
