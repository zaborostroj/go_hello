package kafkaUtils

import (
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewWriter(config Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
}

func CloseWriter(writer *kafka.Writer) {
	err := writer.Close()
	if err != nil {
		log.Print(err)
	}
}

func NewReader(config Config) *kafka.Reader {
	readerConfig := kafka.ReaderConfig{
		Brokers:        config.Brokers,
		GroupID:        config.GroupID,
		Topic:          config.Topic,
		MinBytes:       1,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		Logger:         log.New(os.Stdout, "KAFKA-INFO: ", log.LstdFlags),
		ErrorLogger:    log.New(os.Stderr, "KAFKA-ERR: ", log.LstdFlags),
	}
	return kafka.NewReader(readerConfig)
}

func CloseReader(reader *kafka.Reader) {
	err := reader.Close()
	if err != nil {
		log.Print(err)
	}
}
