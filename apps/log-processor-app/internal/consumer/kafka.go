package consumer

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"worker.go/internal/config"
)

type KafkaConsumer struct {
	reader *kafka.Reader
}

func New(cfg *config.WorkerConfig) *KafkaConsumer {
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.KafkaBrokers,
		GroupID:     "consumer-group-id-01",
		Topic:       cfg.KafkaTopic,
		StartOffset: kafka.LastOffset,
	})

	return &KafkaConsumer{
		reader: kafkaReader,
	}
}

func (k *KafkaConsumer) Read() (kafka.Message, error) {
	message, err := k.reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	return message, err
}

func (k KafkaConsumer) Close() error {
	return k.reader.Close()
}
