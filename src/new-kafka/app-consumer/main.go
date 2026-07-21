/*
References:

	git: https://github.com/segmentio/kafka-go
	doc: https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme
*/
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pseudoelement/new-kafka/common"
	"github.com/segmentio/kafka-go"
)

func main() {
	// testTopicName := common.TOPIC_EVENTS
	// conn0, err := connect(testTopicName, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// conn1, err := connect(testTopicName, 1)
	// if err != nil {
	// 	panic(err)
	// }
	readWithReader("consumers-1", common.TOPIC_PAYMENTS)
	// readWithReader(testTopicName, "consumers-2", 2)
	// readMessages(conn0, 10, 10e3)

	// if err := conn0.Close(); err != nil {
	// 	fmt.Println("failed to close connection:", err)
	// }
	// if err := conn1.Close(); err != nil {
	// 	fmt.Println("failed to close connection:", err)
	// }
	println("Consuming finished.")
}

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		fmt.Println("failed to dial leader")
	}
	return conn, err
}

// NOTE: better use kafka.NewReader, it can specify offset, partition etc.
// Reads all messages in the partition from the start
// Specify a minimum and maximum size in bytes to read (1 char = 1 byte)
func readMessages(conn *kafka.Conn, minSize int, maxSize int) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(minSize, maxSize) //in bytes
	msg := make([]byte, maxSize)              //set the max length of each message
	for {
		msgSize, err := batch.Read(msg)
		if err != nil {
			break
		}
		fmt.Println("[readMessages] msg: ", string(msg[:msgSize]))
	}

	if err := batch.Close(); err != nil { //make sure to close the batch
		fmt.Println("failed to close batch:", err)
	}
}

// Read from the topic using kafka.Reader
// Readers can use consumer groups (but are not required to)
func readWithReader(groupID string, topics ...string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		// Partition: 1,
		GroupID:     groupID,
		GroupTopics: topics,
		// Topic:       topic,
		MaxBytes:    1000, //per message
		StartOffset: kafka.LastOffset,
	})
	log.Println("Reader init.")

	for {
		// ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		// defer cancel()
		msg, err := r.ReadMessage(context.Background())
		// NOTE: FetchMessage doesn't commit read message, you need to r.CommitMessages() manually
		// msg, err := r.FetchMessage(context.Background())
		// err = r.CommitMessages(context.Background(), msg)
		if err != nil {
			fmt.Printf("r.ReadMessage err:%s \n", err.Error())
			break
		}
		fmt.Printf(
			"message at topic/partition/offset %v/%v/%v: %s = %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value),
		)
	}

	if err := r.Close(); err != nil {
		fmt.Println("failed to close reader:", err)
	}
}
