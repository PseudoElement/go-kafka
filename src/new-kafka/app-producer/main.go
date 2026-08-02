package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pseudoelement/new-kafka/common"
	"github.com/segmentio/kafka-go"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	leaderKafkaIP := os.Getenv("LEADER_BROKER_IP")
	if leaderKafkaIP == "" {
		panic("LEADER_BROKER_IP is not set in .env")
	}
	createNewTopic(common.TOPIC_EVENTS, leaderKafkaIP)
	createNewTopic(common.TOPIC_PAYMENTS, leaderKafkaIP)

	// conn0, err := connect(common.TOPIC_EVENTS, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// writeMessages(context.Background(), conn0, 0)
	writeWithWriter(context.Background())
	// conn1, err := connect(common.TOPIC_EVENTS, 1)
	// if err != nil {
	// 	panic(err)
	// }
	// writeMessages(context.Background(), conn1, 1)
}

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"localhost:9092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader")
	}
	return conn, err
}

// OLD WAY
// Writes the messages in the string slice to the topic
func writeMessages(ctx context.Context, conn *kafka.Conn, partition int) {
	var err error
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case val := <-ticker.C:
			msg := fmt.Sprintf("sent log_%d in partition_%d", time.Now().UnixMilli(), partition)
			key := fmt.Sprintf("event-key-%d", val.Second())
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			_, err = conn.WriteMessages(kafka.Message{
				Value:     []byte(msg),
				Key:       []byte(key),
				Partition: partition,
			})
			if err != nil {
				fmt.Println("failed to write messages:", err)
			} else {
				fmt.Printf("success: %s, key: %s\n", msg, key)
			}
		}
	}
}

func writeWithWriter(ctx context.Context) {
	brokerIpsStr := os.Getenv("KAFKA_BROKERS_IPS")
	brokerIps := strings.Split(brokerIpsStr, ",")
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokerIps,
		Balancer:     &kafka.RoundRobin{},
		RequiredAcks: -1,
	})
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case count := <-ticker.C:
			for range 1 {
				go sendMsg(ctx, w, count.Second())
			}
		}
	}
}

func sendMsg(ctx context.Context, w *kafka.Writer, count int) (msg string, key string, err error) {
	var topic string = common.TOPIC_EVENTS
	if count%3 == 0 {
		topic = common.TOPIC_PAYMENTS
	}

	msg = fmt.Sprintf("sent msg_%d to topic %s", time.Now().UnixMilli(), topic)
	key = fmt.Sprintf("event-key-%d", count)

	err = w.WriteMessages(ctx, kafka.Message{
		Value: []byte(msg),
		Key:   []byte(key),
		Topic: topic,
	})
	if err != nil {
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			fmt.Printf("Topic %s not created. Creating...\n", topic)
		} else {
			fmt.Println("failed to write messages:", err)
		}
	} else {
		fmt.Printf("success: %s, key: %s\n", msg, key)
	}

	return msg, key, err
}

func createNewTopic(topic string, kafkaIp string) {
	conn, err := kafka.Dial("tcp", kafkaIp)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     4, // ideally 1 partition per 1 consumer
			ReplicationFactor: 1, // equals to number of broker instances running
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
