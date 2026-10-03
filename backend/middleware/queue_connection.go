package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

// RequireQueueConnection rejects the request with 503 when the AMQP
// connection or channel is missing/closed. Register it on the route
// before the handler.
func RequireQueueConnection(conn *amqp.Connection, ch *amqp.Channel) gin.HandlerFunc {
	return func(c *gin.Context) {
		if conn == nil || ch == nil || conn.IsClosed() || ch.IsClosed() {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "Queue connection unavailable",
			})
			return
		}

		c.Next()
	}
}