package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type AppKafka struct {
	producer  *kafka.Producer
	consumers map[string]chan any
	appCtx    context.Context
}

func NewAppKafka(appCtx context.Context) *AppKafka {
	appKafka := &AppKafka{appCtx: appCtx}
	appKafka.createProducer()
	appKafka.MessagesChan()

	return appKafka
}

func (this *AppKafka) SendMessage(topic string, body any) error {
	value, _ := json.Marshal(body)
	err := this.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	},
		nil, // delivery channel
	)

	return err
}

func (this *AppKafka) MessagesChan() chan kafka.Event {
	go func() {
		for e := range this.producer.Events() {
			select {
			case <-this.appCtx.Done():
				return
			default:
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
					} else {
						log.Println("ev.TopicPartition ==>", ev.TopicPartition)
						log.Println("ev.Key ==>", ev.Key)
						log.Println("ev.Key ==>", string(ev.Value))
					}
				}
			}
		}
	}()

	return this.producer.Events()
}

func (this *AppKafka) createProducer() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		/* without docker */
		// "bootstrap.servers": "host1:9092,host2:9092",
		/* with docker */
		"bootstrap.servers":  "kafka:9092",
		"client.id":          "myProducer",
		"message.timeout.ms": 30000,
		"acks":               "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	p.Close()

	this.producer = p
}
