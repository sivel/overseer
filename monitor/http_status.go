package monitor

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sivel/overseer/status"
)

type HTTPStatusConfig struct {
	Type                 string
	URL                  *url.URL
	Codes                []int
	NotificationInterval time.Duration
	Verify               bool
	Timeout              time.Duration
	Method               string
}

type HTTPStatus struct {
	config *HTTPStatusConfig
	status *status.Status
}

func NewHTTPStatus(conf map[string]interface{}) Monitor {
	var err error
	monitor := new(HTTPStatus)

	var pURL *url.URL
	if urlInterface, ok := conf["url"]; ok {
		pURL, err = url.Parse(urlInterface.(string))
		if err != nil {
			log.Fatalf("Invalid URL: %s", conf["url"].(string))
		} else if !ok {
			log.Fatalf("No URL provided")
		}
	}

	var codes []int = []int{200}
	if codesInterface, ok := conf["codes"]; ok {
		for _, code := range codesInterface.([]interface{}) {
			codes = append(codes, int(code.(float64)))
		}
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

	monitor.config = &HTTPStatusConfig{
		Type:                 conf["type"].(string),
		URL:                  pURL,
		Codes:                codes,
		NotificationInterval: notificationInterval,
		Verify:               verify,
		Timeout:              timeout,
		Method:               method,
	}
	monitor.status = status.NewStatus(status.UNKNOWN, status.UNKNOWN, time.Now(), time.Now(), 0, "")
	return monitor
}

func (m *HTTPStatus) Watch(statusChan chan *status.Status) {
	for {
		m.Check()
		statusChan <- m.status
		time.Sleep(time.Second * 10)
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

func checkChanged(current int, last int, startOfLastStatus time.Time) (bool, time.Time) {
	var start time.Time = startOfLastStatus
	var changed bool = false
	if current != last {
		changed = true
		start = time.Now()
	}
	return changed, start
}

/*func shouldNotify(current int, last int, lastNotification time.Time, notificationInterval time.Duration) (bool, time.Time) {
	var notify bool = false
	var notifyTime time.Time = lastNotification
	if current == last && current == status.UP {
		notify = false
		fmt.Print("1")
	} else if current != last && current == status.UP && last == status.UNKNOWN {
		notify = false
		fmt.Print("2")
	} else if time.Since(lastNotification) > notificationInterval {
		fmt.Print("3")
		notify = true
		notifyTime = time.Now()
	}
	return notify, notifyTime

}*/

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

	m.status = status.NewStatus(current, m.status.Current, start, m.status.LastNotification, duration, message)
}
