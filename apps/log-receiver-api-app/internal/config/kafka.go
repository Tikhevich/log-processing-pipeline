package config

import (
	"fmt"
	"os"
)

type KafkaConfig struct {
	KafkaBrokers []string
	KafkaTopic   string
}

func Load() *KafkaConfig {
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopicName := os.Getenv("KAFKA_LOG_TOPIC_NAME")

	return &KafkaConfig{
		KafkaBrokers: []string{fmt.Sprintf("%v:%v", kafkaHost, kafkaPort)},
		KafkaTopic:   kafkaTopicName,
	}
}
