package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/pseudoelement/new-kafka/common"
	"github.com/segmentio/kafka-go"
)

func main() {
	// conn0, err := connect(common.TOPIC_EVENTS, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// writeMessages(context.Background(), conn0, 0)
	writeWithWriter(context.Background(), 1)
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

func writeWithWriter(ctx context.Context, partition int) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"localhost:9092"},
		Balancer:     &kafka.Hash{},
		RequiredAcks: -1,
	})
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case count := <-ticker.C:
			sendMsg(ctx, w, count.Second())
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
		fmt.Println("failed to write messages:", err)
		if errors.Is(err, kafka.UnknownTopicOrPartition) {
			fmt.Printf("Topic %s not created. Creating...\n", topic)
			createNewTopic(topic)
		}
	} else {
		fmt.Printf("success: %s, key: %s\n", msg, key)
	}

	return msg, key, err
}

func createNewTopic(topic string) {
	conn, err := kafka.Dial("tcp", "localhost:9092")
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
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
