package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	uuid "github.com/satori/go.uuid"
)

var (
	kafkaBrokers = []string{"localhost:9093"}
	KafkaTopic   = "sarama_topic"
	enqueued     int
)

func main() {
	arg := os.Args[1]
	count, _ := strconv.Atoi(os.Args[2])
	if count == 0 {
		count = 10
	}
	if arg == "" {
		panic("Service type need to be send through args")
	}

	forever := make(chan bool)
	producer, err := setupProducer()
	if err != nil {
		panic(err)
	} else {
		log.Println("Kafka AsyncProducer up and running!")
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	produceMessages(arg, count, producer, signals)

	log.Printf("Kafka AsyncProducer finished with %d messages produced.", enqueued)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			forever <- true
		}
	}()
	<-forever
}

// setupProducer will create a AsyncProducer and returns it
func setupProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

// produceMessages will send 'testing 123' to KafkaTopic each second, until receive a os signal to stop e.g. control + c
// by the user in terminal
func produceMessages(t string, count int, producer sarama.AsyncProducer, signals chan os.Signal) {
	for i := 1; i <= count; i++ {
		time.Sleep(time.Second)
		var meta []byte
		myuuid := uuid.NewV4()
		switch strings.ToLower(t) {
		case "email":
			meta = []byte(fmt.Sprintf(`{
			"type": "email", 
			"data": {
				"subject":"Alert from Accurics %d", 
				"from": "info@accurics.com", 
				"text":"Hello world",
				 "to": ["cs-girish.talekar@accurics.com","talekar.g@gmail.com"]
				 }
		 }`, i))
		case "slack":
			meta = []byte(fmt.Sprintf(`{"type": "slack", "data": {"text": "Hey hi %d"}}`, i))
		case "jira":
			meta = []byte(fmt.Sprintf(`{"type": "jira", "data": {"fields": {"issuetype":{"name":"Task"},"project":{"key": "ALT"}, "summary": "Making it work using yaml %d"}}}`, i))
		case "teams":
			meta = []byte(fmt.Sprintf(`{"type": "teams", "data": {"text": "Hey hi %d"}}`, i))
		case "pagerduty":
			meta = []byte(fmt.Sprintf(`{"type": "pagerduty", "data": {"incident":{"type":"incident","title":"The server is on fire %d.","service":{"id":"PIESOML","type":"service_reference"},"priority":{"id":"P53ZZH5","type":"priority_reference"},"urgency":"high","incident_key":"%s","body":{"type":"incident_body","details":"A disk is getting full on this machine. You should investigate what is causing the disk to fill, and ensure that there is an automated process in place for ensuring data is rotated (eg. logs should have logrotate around them). If data is expected to stay on this disk forever, you should start planning to scale up to a larger disk."},"escalation_policy":null}}}`, i, myuuid))
		}

		message := &sarama.ProducerMessage{Topic: KafkaTopic,
			Value: sarama.StringEncoder(meta)}
		select {
		case producer.Input() <- message:
			enqueued++
			log.Println("New Message produced")
		case <-signals:
			producer.AsyncClose() // Trigger a shutdown of the producer.
			return
		}
	}
}
