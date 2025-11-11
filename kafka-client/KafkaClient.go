package kafka_client

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewClient(config Config) *Client {
	return &Client{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(config.Brokers...),
			Topic:        config.Topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        config.Brokers,
			GroupID:        config.GroupID,
			Topic:          config.Topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			CommitInterval: time.Second,
		}),
	}
}

func (client *Client) Close() {
	err := client.writer.Close()
	if err != nil {
		log.Print(err)
	}
	err = client.reader.Close()
	if err != nil {
		log.Print(err)
	}
}

func (client *Client) Writer() *kafka.Writer {
	return client.writer
}

func (client *Client) Reader() *kafka.Reader {
	return client.reader
}
