package main

import (
	"github.com/gin-gonic/gin"
	"server.go/internal/api/handler"
	"server.go/internal/config"
	"server.go/internal/producer"
)

func main() {

	kafkaCfg := config.Load()
	kafkaWriter := producer.GetKafkaWriter(kafkaCfg)
	defer kafkaWriter.Close()

	router := gin.Default()
	router.POST("/upload", func(c *gin.Context) {
		handler.ProcessPostData(c, kafkaWriter)
	})
	router.Run()
}
