package logger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sivel/overseer/status"
	"gopkg.in/olivere/elastic.v3"
)

const settings string = `
{
    "mappings": {
        "%s": {
            "properties": {
                "CheckDuration": {
                    "type": "long"
                },
                "Current": {
                    "type": "integer"
                },
                "Last": {
                    "type": "integer"
                },
                "LastNotification": {
                    "type": "date"
                },
                "Loggers": {
                    "index": "not_analyzed",
                    "type": "string"
                },
                "Message": {
                    "index": "not_analyzed",
                    "type": "string"
                },
                "MonitorName": {
                    "index": "not_analyzed",
                    "type": "string"
                },
                "NotificationInterval": {
                    "type": "long"
                },
                "Notifiers": {
                    "index": "not_analyzed",
                    "type": "string"
                },
                "StartOfCurrentStatus": {
                    "type": "date"
                },
                "Time": {
                    "type": "date"
                }
            }
        }
    }
}
`

type ElasticsearchConfig struct {
	Name     string
	Type     string
	Servers  []string
	Username string
	Password string
	Index    string
	DocType  string `json:"doc_type"`
}

type Elasticsearch struct {
	config *ElasticsearchConfig
	client *elastic.Client
}

func NewElasticsearch(conf []byte, filename string) Logger {
	logger := new(Elasticsearch)
	var config ElasticsearchConfig
	json.Unmarshal(conf, &config)
	logger.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	if len(config.Servers) == 0 {
		log.Fatalf("Servers not provided: %s", filename)
	}

	if config.Index == "" {
		config.Index = "overseer"
	}

	if config.DocType == "" {
		config.DocType = "log"
	}

	var options []elastic.ClientOptionFunc
	options = append(options, elastic.SetURL(config.Servers...))
	if config.Username != "" && config.Password != "" {
		options = append(options, elastic.SetBasicAuth(config.Username, config.Password))
	}

	client, err := elastic.NewSimpleClient(options...)
	if err != nil {
		log.Fatal("Could not connect to Elasticsearch (%s): %s", err, filename)
	}

	exists, err := client.IndexExists(config.Index).Do()
	if err != nil {
		log.Fatal("Could not check if Elasticsearch index exists (%s): %s", err, filename)
	}

	if !exists {
		formattedSettings := fmt.Sprintf(settings, config.DocType)
		_, err := client.CreateIndex(config.Index).BodyString(formattedSettings).Do()
		if err != nil {
			log.Fatal("Could not create Elasticsearch index (%s): %s", err, filename)
		}
	}

	logger.client = client

	return logger
}

func (n *Elasticsearch) Name() string {
	return n.config.Name
}

func (n *Elasticsearch) Type() string {
	return n.config.Type
}

func (n *Elasticsearch) Log(stat *status.Status) {
	_, err := n.client.Index().Index(n.config.Index).Type(n.config.DocType).BodyJson(stat).Do()
	if err != nil {
		log.Printf("Error inserting into Elasticsearch: %s", err)
	}
}
