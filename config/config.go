package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/sivel/overseer/monitor"
	"github.com/sivel/overseer/notifier"
)

type Config struct {
}

type NotifierType struct {
	Type string
}

func getNotifiers(configPath string) []notifier.Notifier {
	notifierPath := filepath.Join(configPath, "notifiers")
	files, err := ioutil.ReadDir(notifierPath)
	if err != nil {
		log.Fatalf("Could not list notifiers configuration directory: %s", err)
	}
	var notifiers []notifier.Notifier
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		var tmp NotifierType
		text, err := ioutil.ReadFile(filepath.Join(notifierPath, f.Name()))
		if err != nil {
			log.Printf("Could not read configuration file: %s", f.Name())
			continue
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			log.Printf("Configuration file not valid JSON: %s", f.Name())
			continue
		}
		notifier, err := notifier.GetNotifier(tmp.Type)
		if err != nil {
			continue
		}
		notifiers = append(notifiers, notifier(text, f.Name()))
	}
	return notifiers
}

func getMonitors(configPath string) []monitor.Monitor {
	monitorPath := filepath.Join(configPath, "monitors")
	files, err := ioutil.ReadDir(monitorPath)
	if err != nil {
		log.Fatalf("Could not list monitors configuration directory: %s", err)
	}
	var monitors []monitor.Monitor
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		var tmp NotifierType
		text, err := ioutil.ReadFile(filepath.Join(monitorPath, f.Name()))
		if err != nil {
			log.Printf("Could not read configuration file: %s", f.Name())
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			log.Printf("Configuration file not valid JSON: %s", f.Name())
		}
		monitor, err := monitor.GetMonitor(tmp.Type)
		if err != nil {
			continue
		}
		monitors = append(monitors, monitor(text))
	}
	return monitors
}

func ParseConfig() ([]monitor.Monitor, []notifier.Notifier) {
	var config Config
	configPath, _ := filepath.Abs("/etc/overseer")
	configFile := filepath.Join(configPath, "overseer.json")
	text, err := ioutil.ReadFile(configFile)
	if err == nil {
		json.Unmarshal(text, &config)
	}
	notifiers := getNotifiers(configPath)
	monitors := getMonitors(configPath)
	return monitors, notifiers
}
