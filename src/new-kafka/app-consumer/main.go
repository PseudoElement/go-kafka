/*
References:

	git: https://github.com/segmentio/kafka-go
	doc: https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme
*/
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pseudoelement/new-kafka/common"
	"github.com/segmentio/kafka-go"
)

func main() {
	testTopicName := common.TOPIC_EVENTS
	conn, _ := connect(testTopicName, 0)

	// readMessages(conn, 10, 10e3)
	readWithReader(testTopicName, "consumer-through-kafka 1")

	if err := conn.Close(); err != nil {
		fmt.Println("failed to close connection:", err)
	}
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

// Reads all messages in the partition from the start
// Specify a minimum and maximum size in bytes to read (1 char = 1 byte)
func readMessages(conn *kafka.Conn, minSize int, maxSize int) {
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	batch := conn.ReadBatch(minSize, maxSize) //in bytes

	msg := make([]byte, 10e3) //set the max length of each message
	for {
		msgSize, err := batch.Read(msg)
		if err != nil {
			break
		}
		fmt.Println(string(msg[:msgSize]))
	}

	if err := batch.Close(); err != nil { //make sure to close the batch
		fmt.Println("failed to close batch:", err)
	}
}

// Read from the topic using kafka.Reader
// Readers can use consumer groups (but are not required to)
func readWithReader(topic string, groupID string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"kafka1:9092", "kafka2:9093"},
		GroupID:     groupID,
		Topic:       topic,
		MaxBytes:    100, //per message
		StartOffset: kafka.LastOffset,
		// more options are available
	})

	ctx, _ := context.WithDeadline(context.Background(),
		time.Now().Add(5*time.Second))
	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}

	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
