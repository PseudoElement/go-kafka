package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type AppKafka struct {
	producer  *kafka.Producer
	consumers map[string]chan KafkaMsg
	appCtx    context.Context
}

func NewAppKafka(appCtx context.Context) *AppKafka {
	appKafka := &AppKafka{appCtx: appCtx, consumers: make(map[string]chan KafkaMsg)}
	appKafka.createProducer()
	appKafka.listen()

	return appKafka
}

type ValueResponse struct {
	Value string `json:"value"`
	Idx   int    `json:"idx"`
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

func (this *AppKafka) AddListener(key string) error {
	_, ok := this.consumers[key]
	if ok {
		return fmt.Errorf("Consumer with name %s already exists.\n", key)
	}

	ch := make(chan KafkaMsg)
	this.consumers[key] = ch

	return nil
}

func (this *AppKafka) ConsumerChan(key string) <-chan KafkaMsg {
	return this.consumers[key]
}

func (this *AppKafka) listen() error {
	go func() {
		for {
			select {
			case <-this.appCtx.Done():
				log.Println("Stopping Kafka event listener due to context cancellation")
				this.producer.Close()
				return
			case e, ok := <-this.producer.Events():
				if !ok {
					log.Println("Kafka producer events channel closed")
					return
				}
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Printf("Kafka failed: %v\n", ev.TopicPartition)
					} else {
						if len(ev.Value) > 0 {
							log.Printf("Topic - %s, message: %v \n", *ev.TopicPartition.Topic, string(ev.Value))
							msg := KafkaMsg{Topic: *ev.TopicPartition.Topic, Value: ev.Value}
							for _, ch := range this.consumers {
								ch <- msg
							}
						}
					}
				}
			}
		}
		// for e := range this.producer.Events() {
		// 	switch ev := e.(type) {
		// 	case *kafka.Message:
		// 		if ev.TopicPartition.Error != nil {
		// 			fmt.Printf("1_listener Delivery failed: %v\n", ev.TopicPartition)
		// 		} else {
		// 			if len(ev.Value) > 0 {
		// 				msg := KafkaMsg{Topic: *ev.TopicPartition.Topic, Value: ev.Value}
		// 				for _, ch := range this.consumers {
		// 					ch <- msg
		// 				}
		// 			}
		// 		}
		// 	}
		// }
	}()

	// Wait for message deliveries before shutting down
	this.producer.Flush(15 * 1000)

	return nil
}

func (this *AppKafka) createProducer() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		/* without docker */
		// "bootstrap.servers": "host1:9092,host2:9092",
		/* with docker */
		"bootstrap.servers":  "kafka:9092",
		"message.timeout.ms": 30000,
	})

	if err != nil {
		log.Fatal("Failed to create producer: ", err)
	}

	this.producer = p
}
