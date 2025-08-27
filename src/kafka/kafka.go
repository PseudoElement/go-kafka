package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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

	err := appKafka.EnsureTopicsExist([]string{"porno", "lgbt", "logger", "github"})
	if err != nil {
		log.Printf("Warning: could not ensure topics exist: %v", err)
	}

	return appKafka
}

type ValueResponse struct {
	Value string `json:"value"`
	Idx   int    `json:"idx"`
}

func (this *AppKafka) SendMessage(topic string, body any) error {
	if this.producer.IsClosed() {
		log.Println("Kafka producer is closed now!")
		return fmt.Errorf("Kafka producer is closed now!")
	}

	value, _ := json.Marshal(body)
	err := this.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	},
		nil, // delivery channel
	)

	return err
}

// correct listening
func (this *AppKafka) CreateConsumer(groupID string, topics []string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        "kafka:9092",
		"group.id":                 groupID, // Different group for different modules
		"auto.offset.reset":        "earliest",
		"allow.auto.create.topics": true,
	})

	if err != nil {
		log.Fatal("Failed to create consumer: ", err)
	}

	c.SubscribeTopics(topics, nil)

	return c
}

func (this *AppKafka) ListenViaConsumer(c *kafka.Consumer, callback func(msg KafkaMsg)) {
	for {
		select {
		case <-this.appCtx.Done():
			log.Println("[ListenViaConsumer] context called appCtx.Done()")
			c.Close()
			return
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					log.Printf("Consumer error: %v", err)
				}
				continue
			}

			kafkaMsg := KafkaMsg{
				Topic: *msg.TopicPartition.Topic,
				Value: msg.Value,
			}

			callback(kafkaMsg)
		}
	}
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

// listen via channel (self-written)
func (this *AppKafka) ConsumerChan(key string) <-chan KafkaMsg {
	return this.consumers[key]
}

// !!! IN KAFKA YOU NEED MANUALLY CREATE TOPICS
func (this *AppKafka) EnsureTopicsExist(topics []string) error {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
	})
	if err != nil {
		return err
	}
	defer admin.Close()

	var topicSpecs []kafka.TopicSpecification
	for _, topic := range topics {
		topicSpecs = append(topicSpecs, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1, // Adjust as needed
			ReplicationFactor: 1, // Adjust for your cluster
		})
	}

	// Try to create topics (they might already exist, which is fine)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = admin.CreateTopics(ctx, topicSpecs)
	// It's OK if topics already exist, we ignore that error
	if err != nil && !isTopicExistsError(err) {
		return err
	}

	return nil
}

func isTopicExistsError(err error) bool {
	if kafkaErr, ok := err.(kafka.Error); ok {
		return kafkaErr.Code() == kafka.ErrTopicAlreadyExists
	}
	return false
}

func (this *AppKafka) listen() error {
	go func() {
		for {
			select {
			case <-this.appCtx.Done():
				log.Println("[listen] triggerer appCtx.Done()")
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
