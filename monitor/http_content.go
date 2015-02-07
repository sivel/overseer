package monitor

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/sivel/overseer/status"
)

type HTTPContentConfig struct {
	Name                 string
	URL                  *url.URL
	Content              string
	CheckInterval        time.Duration
	NotificationInterval time.Duration
	Verify               bool
	Timeout              time.Duration
	Method               string
}

type HTTPContent struct {
	config *HTTPContentConfig
	status *status.Status
}

func NewHTTPContent(conf map[string]interface{}) Monitor {
	var err error
	monitor := new(HTTPContent)

	var pURL *url.URL
	if urlInterface, ok := conf["url"]; ok {
		pURL, err = url.Parse(urlInterface.(string))
		if err != nil {
			log.Fatalf("Invalid URL: %s", conf["url"].(string))
		} else if !ok {
			log.Fatalf("No URL provided")
		}
	}

	var content string
	if contentInterface, ok := conf["content"]; ok {
		content = contentInterface.(string)
	} else {
		log.Fatal("No content match provided")
	}

	var name string = pURL.String()
	if nameInterface, ok := conf["name"]; ok {
		name = nameInterface.(string)
	}

	var checkInterval time.Duration = time.Second * 10
	if ci, ok := conf["check_interval"]; ok {
		checkInterval, err = time.ParseDuration(ci.(string))
	}

	var notificationInterval time.Duration = time.Second * 60
	if ni, ok := conf["notification_interval"]; ok {
		notificationInterval, err = time.ParseDuration(ni.(string))
	}

	var verify bool = false
	if verifyInterface, ok := conf["verify"]; ok {
		verify = verifyInterface.(bool)
	}

	var timeout time.Duration = time.Second * 2
	if timeoutInterface, ok := conf["timeout"]; ok {
		timeout, _ = time.ParseDuration(timeoutInterface.(string))
	}

	var method string = "HEAD"
	if methodInterface, ok := conf["method"]; ok {
		method = strings.ToUpper(methodInterface.(string))
	}

	monitor.config = &HTTPContentConfig{
		Name:                 name,
		URL:                  pURL,
		Content:              content,
		CheckInterval:        checkInterval,
		NotificationInterval: notificationInterval,
		Verify:               verify,
		Timeout:              timeout,
		Method:               method,
	}
	monitor.status = status.NewStatus(
		name,
		status.UNKNOWN,
		status.UNKNOWN,
		notificationInterval,
		time.Now(),
		time.Now(),
		0,
		"",
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
	fmt.Println("HTTPContent Check Running for " + m.config.URL.String())

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
	resp, err := client.Do(&http.Request{Method: "GET", URL: m.config.URL})
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
	)
}
