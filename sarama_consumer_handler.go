package main

import (
	"log"

	"github.com/Shopify/sarama"
	"github.com/girishg4t/go-notifier/pkg/libs"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ConsumerGroupHandler represents the sarama consumer group
type ConsumerGroupHandler struct{}

// Setup is run before consumer start consuming, is normally used to setup things such as database connections
func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages(), here is supposed to be what you want to
// do with the message. In this example the message will be logged with the topic name, partition and message value.
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		//fmt.Printf("Message topic:%q partition:%d offset:%d message: %v\n",msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		handleMessage(msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func handleMessage(b []byte) {
	//Parse the incoming message from kafka
	var msg libs.Message
	if err := json.Unmarshal(b, &msg); err != nil {
		log.Panicf("Not able to parse the message %s", string(b))
		return
	}

	//Bind message to template with data
	ys, err := libs.GetParsedTemplate("./templates/jira.yaml", msg.Data)
	if err != nil || ys == "" {
		log.Panicf("Not able to parse the template %s", ys)
		return
	}

	//Convert yaml string to configuration object
	var c libs.Configuration
	err = yaml.Unmarshal([]byte(ys), &c)
	if err != nil {
		log.Fatalf("error while converting to configuration struct: %v", err)
		return
	}
	log.Println(c)
}
