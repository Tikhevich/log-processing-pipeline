package producer

import (
	"github.com/segmentio/kafka-go"
	"server.go/internal/config"
)

func GetKafkaWriter(cfg *config.KafkaConfig) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(cfg.KafkaBrokers...),
		Topic:                  cfg.KafkaTopic,
		Balancer:               &kafka.LeastBytes{},
		BatchSize:              1,
		AllowAutoTopicCreation: true,
	}
}
