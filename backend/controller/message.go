// controller/message_controller.go

package controller

import (
	"net/http"

	"backend/dto"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageController struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

func (mc *MessageController) CreateMessage(c *gin.Context) {
	var request dto.CreateMessageRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	err := mc.Channel.Publish(
		"",
		mc.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(request.Message),
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish message",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Message sent to RabbitMQ",
	})
}