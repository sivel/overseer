package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/sivel/overseer/logger"
	"github.com/sivel/overseer/monitor"
	"github.com/sivel/overseer/notifier"
)

type Config struct {
}

type PluginType struct {
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
		var tmp PluginType
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
			log.Printf("%s: %s", err.Error(), f.Name())
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
		var tmp PluginType
		text, err := ioutil.ReadFile(filepath.Join(monitorPath, f.Name()))
		if err != nil {
			log.Printf("Could not read configuration file: %s", f.Name())
			continue
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			log.Printf("Configuration file not valid JSON: %s", f.Name())
			continue
		}
		monitor, err := monitor.GetMonitor(tmp.Type)
		if err != nil {
			log.Printf("%s: %s", err.Error(), f.Name())
			continue
		}
		monitors = append(monitors, monitor(text, f.Name()))
	}
	return monitors
}

func getLoggers(configPath string) []logger.Logger {
	loggerPath := filepath.Join(configPath, "loggers")
	files, err := ioutil.ReadDir(loggerPath)
	if err != nil {
		log.Fatalf("Could not list loggers configuration directory: %s", err)
	}
	var loggers []logger.Logger
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		var tmp PluginType
		text, err := ioutil.ReadFile(filepath.Join(loggerPath, f.Name()))
		if err != nil {
			log.Printf("Could not read configuration file: %s", f.Name())
			continue
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			log.Printf("Configuration file not valid JSON: %s", f.Name())
			continue
		}
		logger, err := logger.GetLogger(tmp.Type)
		if err != nil {
			log.Printf("%s: %s", err.Error(), f.Name())
			continue
		}
		loggers = append(loggers, logger(text, f.Name()))
	}
	return loggers
}

func ParseConfig() ([]monitor.Monitor, []notifier.Notifier, []logger.Logger) {
	configPath := "/etc/overseer"
	notifiers := getNotifiers(configPath)
	monitors := getMonitors(configPath)
	loggers := getLoggers(configPath)
	return monitors, notifiers, loggers
}
