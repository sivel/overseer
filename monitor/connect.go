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
	CheckIntervalString        string `json:"check_interval"`
	CheckInterval              time.Duration
	NotificationIntervalString string `json:"notification_interval"`
	NotificationInterval       time.Duration
	TimeoutString              string `json:"timeout"`
	Timeout                    time.Duration
	Notifiers                  []string
}

type Connect struct {
	config *ConnectConfig
	status *status.Status
}

func NewConnect(conf []byte) Monitor {
	var err error
	monitor := new(Connect)
	var config ConnectConfig
	err = json.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %s", err)
	} else {
		monitor.config = &config
	}

	if config.Host == "" {
		log.Fatalf("No Host provided")
	}

	if config.Port == 0 {
		log.Fatalf("No Port provided")
	}

	if config.Name == "" {
		config.Name = fmt.Sprintf("%s:%d", config.Host, config.Port)
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

func (m *Connect) Watch(statusChan chan *status.Status) {
	for {
		m.Check()
		statusChan <- m.status
		time.Sleep(m.config.CheckInterval)
	}
}

func (m *Connect) Check() {
	requestStart := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), m.config.Timeout)
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
	)
}
