package monitor

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/sivel/overseer/status"
)

type HTTPStatusConfig struct {
	Name                       string
	URLString                  string `json:"url"`
	URL                        *url.URL
	Codes                      []int
	CheckIntervalString        string `json:"check_interval"`
	CheckInterval              time.Duration
	NotificationIntervalString string `json:"notification_interval"`
	NotificationInterval       time.Duration
	Verify                     bool
	TimeoutString              string `json:"timeout"`
	Timeout                    time.Duration
	Method                     string
	Notifiers                  []string
	Loggers                    []string
	Retries                    int
}

type HTTPStatus struct {
	config *HTTPStatusConfig
	status *status.Status
}

func NewHTTPStatus(conf []byte, filename string) Monitor {
	var err error
	monitor := new(HTTPStatus)
	var config HTTPStatusConfig
	json.Unmarshal(conf, &config)
	monitor.config = &config

	if config.URLString == "" {
		log.Fatalf("No URL provided: %s", filename)
	}

	config.URL, err = url.Parse(config.URLString)
	if err != nil {
		log.Fatalf("Invalid URL (%s) provided: %s", config.URLString, filename)
	}

	if config.Name == "" {
		config.Name = config.URL.String()
	}

	if len(config.Codes) == 0 {
		config.Codes = []int{200}
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
		config.Method = "HEAD"
	}

	if config.Retries == 0 {
		config.Retries = 3
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
		[]string{},
	)
	return monitor
}

func (m *HTTPStatus) Watch(statusChan chan *status.Status) {
	for {
		m.Check()
		statusChan <- m.status
		time.Sleep(m.config.CheckInterval)
	}
}

func isValidCode(code int, codes []int) bool {
	var valid bool = false
	for _, c := range codes {
		if c == code {
			valid = true
		}
	}
	return valid
}

func (m *HTTPStatus) check() (int, string) {
	transport := http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, m.config.Timeout)
		},
	}

	if m.config.URL.Scheme == "https" {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: !m.config.Verify}
	}

	client := http.Client{
		Transport: &transport,
	}

	resp, err := client.Do(&http.Request{Method: m.config.Method, URL: m.config.URL})

	var current int = status.UP
	var message string = "OK"
	if err != nil {
		current = status.DOWN
		message = err.Error()
	} else {
		defer resp.Body.Close()

		if !isValidCode(resp.StatusCode, m.config.Codes) {
			current = status.DOWN
			message = fmt.Sprintf("Invalid response code: %d", resp.StatusCode)
		}
	}

	return current, message
}

func (m *HTTPStatus) Check() {
	var current int
	var message string

	requestStart := time.Now()
	for i := 0; i < m.config.Retries; i++ {
		current, message = m.check()
		if current == status.UP {
			break
		}
		time.Sleep(m.config.Timeout)
	}
	duration := time.Now().UnixNano() - requestStart.UnixNano()

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
		m.config.Loggers,
	)
}
