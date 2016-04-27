package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sivel/overseer/status"
)

type SlackConfig struct {
	Name       string
	Type       string
	WebhookURL string `json:"webhook_url"`
	Channel    string
	Username   string
}

type Slack struct {
	config *SlackConfig
}

type Message struct {
	Text        string        `json:"text"`
	Username    string        `json:"username"`
	IconUrl     string        `json:"icon_url"`
	IconEmoji   string        `json:"icon_emoji"`
	Channel     string        `json:"channel"`
	UnfurlLinks bool          `json:"unfurl_links"`
	Attachments []*Attachment `json:"attachments"`
}

func (m *Message) NewAttachment() *Attachment {
	a := &Attachment{}
	m.AddAttachment(a)
	return a
}

func (m *Message) AddAttachment(a *Attachment) {
	m.Attachments = append(m.Attachments, a)
}

type Attachment struct {
	Fallback string   `json:"fallback"`
	Text     string   `json:"text"`
	Pretext  string   `json:"pretext"`
	Color    string   `json:"color"`
	Fields   []*Field `json:"fields"`
	MrkdwnIn []string `json:"mrkdwn_in"`
}

func (a *Attachment) NewField() *Field {
	f := &Field{}
	a.AddField(f)
	return f
}

func (a *Attachment) AddField(f *Field) {
	a.Fields = append(a.Fields, f)
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func NewSlack(conf []byte, filename string) Notifier {
	notifier := new(Slack)

	var config SlackConfig
	json.Unmarshal(conf, &config)
	notifier.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	if config.WebhookURL == "" {
		log.Fatalf("Slack webhook URL not provided: %s", filename)
	}

	return notifier
}

func (n *Slack) Name() string {
	return n.config.Name
}

func (n *Slack) Type() string {
	return n.config.Type
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

	msg := Message{
		Text: fmt.Sprintf("[%s] %s", stateString(stat), stat.MonitorName),
	}

	if n.config.Username != "" {
		msg.Username = n.config.Username
	}

	if n.config.Channel != "" {
		msg.Channel = n.config.Channel
	}

	attach := msg.NewAttachment()
	attach.Color = slackColor(stat)

	field := attach.NewField()
	field.Title = stateString(stat)
	field.Value = fmt.Sprintf(
		"%s [%dms] [%s]",
		stat.Message,
		stat.CheckDuration/1000000,
		time.Since(stat.StartOfCurrentStatus),
	)

	body, _ := json.Marshal(msg)
	buf := bytes.NewReader(body)

	resp, err := http.Post(n.config.WebhookURL, "application/json", buf)
	if err != nil {
		log.Print("Slack notifier: unable to send message")
		return
	}

	if resp.StatusCode != 200 {
		log.Print("Slack notifier: unable to send message")
		return
	}
}
