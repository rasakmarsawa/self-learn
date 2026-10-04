package main

import (
	"log"

	"auth/controller"
	"auth/middleware"
	"auth/connection"

	"github.com/gin-gonic/gin"
)

func main() {
	redisClient := connection.InitRedis()

	authController := controller.NewAuthController(redisClient)

	router := gin.Default()

	router.GET("/hello",
		authController.Hello,
	)

	router.POST("/login",
		authController.Login,
	)

	router.POST("/refresh",
		authController.Refresh,
	)

	router.POST("/logout",
		authController.Logout,
	)	

	router.GET("/authenticated-hello",
		middleware.ValidateToken(),
		authController.Hello,
	)

	log.Println("Auth running on :8081")

	if err := router.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}