package main

import (
	"log"

	"auth/controller"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello",
		controller.Hello,
	)

	router.POST("/login",
		controller.Login,
	)

	router.POST("/authenticated-hello",
		middleware.ValidateToken(),
		controller.Hello,
	)

	log.Println("Backend running on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}