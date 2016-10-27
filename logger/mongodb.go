package logger

import (
	"encoding/json"
	"log"

	"labix.org/v2/mgo"

	"github.com/sivel/overseer/status"
)

type MongoDBConfig struct {
	Name       string
	Type       string
	MongoDBUri string `json:"mongodb_uri"`
}

type MongoDB struct {
	config  *MongoDBConfig
	session *mgo.Session
}

func NewMongoDB(conf []byte, filename string) Logger {
	logger := new(MongoDB)
	var config MongoDBConfig
	json.Unmarshal(conf, &config)
	logger.config = &config

	if config.Name == "" {
		config.Name = config.Type
	}

	if config.MongoDBUri == "" {
		log.Fatalf("MongoDBUri not provided: %s", filename)
	}

	mongo, err := mgo.Dial(config.MongoDBUri)
	if err != nil {
		log.Fatal("Could not connect to MongoDB (%s): %s", err, filename)
	}

	if mongo.DB("").Name == "test" {
		log.Fatalf("The provided Mongo Connection String URI does not appear to have a database name: %s", filename)
	}

	logger.session = mongo

	return logger
}

func (n *MongoDB) Name() string {
	return n.config.Name
}

func (n *MongoDB) Type() string {
	return n.config.Type
}

func (n *MongoDB) Log(stat *status.Status) {
	mongo := n.session.Copy()
	defer mongo.Close()
	c := mongo.DB("").C("logs")

	err := c.Insert(stat)
	if err != nil {
		log.Printf("Error inserting into MongoDB: %s", err)
	}

}
