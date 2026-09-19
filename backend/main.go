package main

import (
	"log"

	"backend/connection"
	"backend/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	conn, ch, queue, err := connection.InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	defer ch.Close()

	router := gin.Default()

	messageController := controller.MessageController{
		Channel: ch,
		Queue:   queue,
	}

	router.POST("/messages", messageController.CreateMessage)

	log.Println("Backend running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}