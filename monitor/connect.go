package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/sivel/overseer/status"
)

type ConnectConfig struct {
	Name                       string
	Host                       string
	Port                       int
	Protocol                   string
	CheckIntervalString        string `json:"check_interval"`
	CheckInterval              time.Duration
	NotificationIntervalString string `json:"notification_interval"`
	NotificationInterval       time.Duration
	TimeoutString              string `json:"timeout"`
	Timeout                    time.Duration
	Notifiers                  []string
	Loggers                    []string
	Retries                    int
}

type Connect struct {
	config *ConnectConfig
	status *status.Status
}

func NewConnect(conf []byte, filename string) Monitor {
	var err error
	monitor := new(Connect)
	var config ConnectConfig
	json.Unmarshal(conf, &config)
	monitor.config = &config

	if config.Host == "" {
		log.Fatalf("No Host provided: %s", filename)
	}

	if config.Port == 0 {
		log.Fatalf("No Port provided: %s", filename)
	}

	if config.Name == "" {
		config.Name = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	if config.Protocol == "" {
		config.Protocol = "tcp"
	}

	if config.Protocol != "tcp" || config.Protocol != "udp" {
		log.Fatalf("Invalid Protocol specified: %s", filename)
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

func (m *Connect) Watch(statusChan chan *status.Status) {
	for {
		m.Check()
		statusChan <- m.status
		time.Sleep(m.config.CheckInterval)
	}
}

func (m *Connect) Check() {
	var err error = nil
	var conn net.Conn

	requestStart := time.Now()
	for i := 0; i < m.config.Retries; i++ {
		conn, err = net.DialTimeout(m.config.Protocol, fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), m.config.Timeout)
		if err == nil {
			break
		}
		time.Sleep(m.config.Timeout)
	}
	duration := time.Now().UnixNano() - requestStart.UnixNano()

	var message string = "OK"
	var current int = status.UNKNOWN
	if err != nil {
		current = status.DOWN
		message = err.Error()
	} else {
		conn.Close()
		current = status.UP
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
		m.config.Loggers,
	)
}
