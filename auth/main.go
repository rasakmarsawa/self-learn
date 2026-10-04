package main

import (
	"log"

	"auth/controller"
	"auth/middleware"
	"auth/connection"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	redisClient := connection.InitRedis()

	authController := controller.NewAuthController(redisClient)

	router := gin.Default()


	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))	

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