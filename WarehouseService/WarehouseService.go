package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/kafka-client"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Настройки Kafka
	cfg := kafka_client.Config{
		Brokers: []string{"localhost:29092"},
		Topic:   "orders",
		GroupID: "warehouse-group",
	}

	client := kafka_client.NewClient(cfg)
	defer client.Close()

	reader := client.Reader()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch SIGINT/SIGTERM for the correct shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received shutdown signal, exiting...")
		cancel()
	}()

	log.Println("Kafka listener started...")

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("🛑 Контекст отменён — слушатель завершает работу")
				break
			}
			if errors.Is(err, kafka.ErrGenerationEnded) {
				log.Println("🔄 Rebalance detected — consumer generation ended")
				continue
			}
			if err.Error() == "EOF" {
				log.Println("📭 EOF от Kafka — ожидание новых сообщений...")
				time.Sleep(time.Second)
				continue
			}
			log.Printf("❌ Ошибка чтения сообщения: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Received message from topic %s: key=%s, value=%s, offset=%d",
			m.Topic, string(m.Key), string(m.Value), m.Offset)

		if err := reader.CommitMessages(ctx, m); err != nil {
			log.Printf("⚠️ Ошибка коммита offset’а: %v", err)
		} else {
			log.Printf("✅ Offset зафиксирован: partition=%d offset=%d", m.Partition, m.Offset)
		}
	}
}
