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
	logger.appKafka.AddListener("logger")
	go logger.listenChan()

	return logger
}

func (this *LoggerModule) listenChan() {
	for msg := range this.appKafka.ConsumerChan("logger") {
		if msg.Topic == "logger_topic" {
			log.Println("LoggerModule message - ", string(msg.Value))
		}
	}
}
