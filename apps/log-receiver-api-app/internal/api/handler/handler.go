package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func ProcessPostData(c *gin.Context, writer *kafka.Writer) {

	data, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	messages := make([]kafka.Message, 0, 1)
	messages = append(messages, kafka.Message{
		Value: data,
	})

	err = writer.WriteMessages(c.Request.Context(), messages...)

	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logs processed",
	})
}
