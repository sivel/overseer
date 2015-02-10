package monitor

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/sivel/overseer/status"
)

type HTTPContentConfig struct {
	Name                       string
	URLString                  string `json:"url"`
	URL                        *url.URL
	Content                    string
	CheckIntervalString        string `json:"check_interval"`
	CheckInterval              time.Duration
	NotificationIntervalString string `json:"notification_interval"`
	NotificationInterval       time.Duration
	Verify                     bool
	TimeoutString              string `json:"timeout"`
	Timeout                    time.Duration
	Method                     string
	Notifiers                  []string
}

type HTTPContent struct {
	config *HTTPContentConfig
	status *status.Status
}

func NewHTTPContent(conf []byte, filename string) Monitor {
	var err error
	monitor := new(HTTPContent)
	var config HTTPContentConfig
	json.Unmarshal(conf, &config)
	monitor.config = &config

	if config.URLString == "" {
		log.Fatalf("No URL provided: %s", filename)
	}

	config.URL, err = url.Parse(config.URLString)
	if err != nil {
		log.Fatalf("Invalid URL (%s) provided: %s", config.URLString, filename)
	}

	if config.Content == "" {
		log.Fatalf("No content match provided: %s", filename)
	}

	if config.Name == "" {
		config.Name = config.URL.String()
	}

	config.CheckInterval, err = time.ParseDuration(config.CheckIntervalString)
	if err != nil {
		config.CheckInterval = time.Second * 10
	}

	config.NotificationInterval, err = time.ParseDuration(config.NotificationIntervalString)
	if err != nil {
		config.NotificationInterval = time.Second * 60
	}

	config.Timeout, err = time.ParseDuration(config.TimeoutString)
	if err != nil {
		config.Timeout = time.Second * 2
	}

	if config.Method == "" {
		config.Method = "GET"
	}

	monitor.status = status.NewStatus(
		config.Name,
		status.UNKNOWN,
		status.UNKNOWN,
		config.NotificationInterval,
		time.Now(),
		time.Now(),
		0,
		"",
		[]string{},
	)
	return monitor
}

func (m *HTTPContent) Watch(statusChan chan *status.Status) {
	for {
		m.Check()
		statusChan <- m.status
		time.Sleep(m.config.CheckInterval)
	}
}

func (m *HTTPContent) Check() {
	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, m.config.Timeout)
		},
	}

	if m.config.URL.Scheme == "https" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: m.config.Verify}
	}

	client := http.Client{
		Transport: &transport,
	}

	requestStart := time.Now()
	resp, err := client.Do(&http.Request{Method: m.config.Method, URL: m.config.URL})
	duration := time.Now().UnixNano() - requestStart.UnixNano()

	var current int = status.UP
	var message string = "OK"
	if err != nil {
		current = status.DOWN
		message = err.Error()
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			current = status.DOWN
			message = fmt.Sprintf("Could not retrieve body: %s", err)
		} else {
			re := regexp.MustCompile(m.config.Content)
			if !re.Match(body) {
				current = status.DOWN
				message = fmt.Sprintf("No content match")
			}
		}
	}

	_, start := checkChanged(current, m.status.Current, m.status.StartOfCurrentStatus)

	m.status = status.NewStatus(
		m.config.Name,
		current,
		m.status.Current,
		m.config.NotificationInterval,
		start,
		m.status.LastNotification,
		duration,
		message,
		m.config.Notifiers,
	)
}
