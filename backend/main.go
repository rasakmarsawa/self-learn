package main

import (
	"log"

	"backend/connection"
	"backend/controller"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	conn, ch, queue, err := connection.InitRabbitMQ()
	if err != nil {
		log.Printf("rabbitmq unavailable, message endpoint disabled: %v", err)
	}

	var messageController controller.MessageController
	if err == nil {
		defer conn.Close()
		defer ch.Close()
		messageController = controller.MessageController{
			Channel: ch,
			Queue:   queue,
		}
	}

	router := gin.Default()

	router.POST("/messages",
		middleware.RequireQueueConnection(conn, ch),
		messageController.CreateMessage,
	)

	router.POST("/ask", 
		messageController.Ask,
	)

	log.Println("Backend running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}