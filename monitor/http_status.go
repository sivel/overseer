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
}

type HTTPStatus struct {
	config *HTTPStatusConfig
	status *status.Status
}

func NewHTTPStatus(conf []byte) Monitor {
	var err error
	monitor := new(HTTPStatus)
	var config HTTPStatusConfig
	err = json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	} else {
		monitor.config = &config
	}

	if config.URLString == "" {
		log.Fatalf("No URL provided")
	}

	config.URL, err = url.Parse(config.URLString)
	if err != nil {
		log.Fatalf("Invalid URL provided: %s", config.URLString)
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

	monitor.status = status.NewStatus(
		config.Name,
		status.UNKNOWN,
		status.UNKNOWN,
		config.NotificationInterval,
		time.Now(),
		time.Now(),
		0,
		"",
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

func (m *HTTPStatus) Check() {
	fmt.Println("HTTPStatus Check Running for " + m.config.URL.String())

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
	resp, err := client.Do(&http.Request{Method: "HEAD", URL: m.config.URL})
	duration := time.Now().UnixNano() - requestStart.UnixNano()

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
	)
}
