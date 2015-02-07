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

func getNotifiers(configPath string) []notifier.Notifier {
	notifierPath := filepath.Join(configPath, "notifiers")
	files, err := ioutil.ReadDir(notifierPath)
	if err != nil {
		log.Fatalf("%s", err)
	}
	var notifiers []notifier.Notifier
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		var tmp interface{}
		text, err := ioutil.ReadFile(filepath.Join(notifierPath, f.Name()))
		if err != nil {
			continue
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			continue
		}
		nConf := tmp.(map[string]interface{})
		if notifierType, ok := nConf["type"]; ok {
			nc, err := notifier.GetNotifier(notifierType.(string))
			if err != nil {
				continue
			}
			conf := nc(nConf)
			notifiers = append(notifiers, conf)
		}
	}
	return notifiers
}

func getMonitors(configPath string) []monitor.Monitor {
	monitorPath := filepath.Join(configPath, "monitors")
	files, err := ioutil.ReadDir(monitorPath)
	if err != nil {
		log.Fatalf("%s", err)
	}
	var monitors []monitor.Monitor
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		var tmp interface{}
		text, err := ioutil.ReadFile(filepath.Join(monitorPath, f.Name()))
		if err != nil {
			continue
		}
		err = json.Unmarshal(text, &tmp)
		if err != nil {
			continue
		}
		mConf := tmp.(map[string]interface{})
		if monitorType, ok := mConf["type"]; ok {
			mc, err := monitor.GetMonitor(monitorType.(string))
			if err != nil {
				continue
			}
			conf := mc(mConf)
			monitors = append(monitors, conf)
		}
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
