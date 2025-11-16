package warehouse

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/zaborostroj/go_hello/internal"
	config "example.com/zaborostroj/go_hello/pkg/shared"
	"github.com/segmentio/kafka-go"
)

func Start() {

	appCfg := config.LoadConfig[ServiceConfig]("warehouse-service")

	cfg := internal.Config{
		Brokers: []string{appCfg.KAFKA.Host + ":" + appCfg.KAFKA.Port},
		Topic:   appCfg.KAFKA.Topic,
		GroupID: appCfg.KAFKA.GroupId,
	}
	kafkaClient := internal.NewReader(cfg)
	defer internal.CloseReader(kafkaClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Catch SIGINT/SIGTERM for proper shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Received shutdown signal, exiting...")
		cancel()
	}()

	log.Println("Kafka listener started...")

	for {
		m, err := kafkaClient.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("🛑 Context cancelled — listener is shutting down")
				break
			}
			if errors.Is(err, kafka.ErrGenerationEnded) {
				log.Println("🔄 Rebalance detected — consumer generation ended")
				continue
			}
			if err.Error() == "EOF" {
				log.Println("📭 EOF from Kafka — waiting for new messages...")
				time.Sleep(time.Second)
				continue
			}
			log.Printf("❌ Error reading message: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Received message from topic %s: key=%s, value=%s, offset=%d",
			m.Topic, string(m.Key), string(m.Value), m.Offset)

		if err := kafkaClient.CommitMessages(ctx, m); err != nil {
			log.Printf("⚠️ Error committing offset: %v", err)
		} else {
			log.Printf("✅ Offset committed: partition=%d offset=%d", m.Partition, m.Offset)
		}
	}
}
