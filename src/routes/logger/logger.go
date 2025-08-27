package logger

import (
	"log"

	"github.com/pseudoelement/go-kafka/src/kafka"
)

type LoggerModule struct {
	appKafka *kafka.AppKafka
}

func NewLoggerModule(appKafka *kafka.AppKafka) *LoggerModule {
	logger := &LoggerModule{appKafka: appKafka}

	// listen via channel
	logger.appKafka.AddListener("logger")
	go logger.listenChan()

	// listen via *kafka.Consumer
	consumer := appKafka.CreateConsumer("1", []string{"porno", "lgbt"})
	go appKafka.ListenViaConsumer(consumer, func(msg kafka.KafkaMsg) {
		log.Printf("[appKafka.ListenViaConsumer] value - %v, topic - %s\n", string(msg.Value), msg.Topic)
	})

	// listen via *kafka.Consumer
	consumer2 := appKafka.CreateConsumer("2", []string{"porno", "lgbt"})
	go appKafka.ListenViaConsumer(consumer2, func(msg kafka.KafkaMsg) {
		log.Printf("[appKafka.ListenViaConsumer_2] value - %v, topic - %s\n", string(msg.Value), msg.Topic)
	})

	return logger
}

func (this *LoggerModule) listenChan() {
	for msg := range this.appKafka.ConsumerChan("logger") {
		if msg.Topic == "logger_topic" {
			log.Println("LoggerModule message - ", string(msg.Value))
		}
	}
}
