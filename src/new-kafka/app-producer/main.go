package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// connect("events", TOPIC_EVENTS)
	// writeMessages()
}

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp",
		"kafka1:9092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader")
	}
	return conn, err
}

// Writes the messages in the string slice to the topic
func writeMessages(ctx context.Context, conn *kafka.Conn) {
	var err error
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case val := <-ticker.C:
			msg := fmt.Sprintf("[app_producer] sent log_%d", val.Second())
			_, err = conn.WriteMessages(kafka.Message{
				Value: []byte(msg),
				Topic: "events",
			})
			if err != nil {
				fmt.Println("failed to write messages:", err)
			}
		}
	}
}
