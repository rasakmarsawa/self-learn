package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"backend/dto"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/joho/godotenv"
)

type MessageController struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

var (
	chatEndpoint string
	chatAPIKey   string
	chatModel    string
	chatClient   = &http.Client{Timeout: 0}
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	chatEndpoint = os.Getenv("CHAT_ENDPOINT")
	chatAPIKey = os.Getenv("CHAT_API_KEY")
	chatModel = os.Getenv("CHAT_MODEL")

	if chatModel == "" {
		chatModel = "Atria-Dawn-Preview"
	}
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

// Ask proxies the request to the upstream chat API with stream=true,
// logs every raw SSE line as it arrives, and returns the accumulated
// reasoning and text as JSON.
func (mc *MessageController) Ask(c *gin.Context) {
	var request dto.ChatRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if len(request.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "messages is required",
		})
		return
	}

	if request.Model == "" {
		request.Model = chatModel
	}

	request.Stream = true

	payload, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to encode upstream request",
		})
		return
	}

	log.Printf("[ask] -> upstream %s body=%s", chatEndpoint, payload)

	upstream, err := http.NewRequestWithContext(
		c.Request.Context(),
		http.MethodPost,
		chatEndpoint,
		bytes.NewReader(payload),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to build upstream request",
		})
		return
	}

	log.Printf("[key] %s", chatAPIKey)

	upstream.Header.Set("Authorization", "Bearer "+chatAPIKey)
	upstream.Header.Set("Content-Type", "application/json")
	upstream.Header.Set("Accept", "text/event-stream")
	upstream.Header.Set("Cache-Control", "no-cache")

	response, err := chatClient.Do(upstream)
	if err != nil {
		log.Printf("[ask] upstream transport error: %v", err)

		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Failed to reach upstream",
		})
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(io.LimitReader(response.Body, 4096))

		log.Printf(
			"[ask] upstream status %d body=%s",
			response.StatusCode,
			raw,
		)

		c.JSON(http.StatusBadGateway, gin.H{
			"error":  "Upstream returned an error",
			"status": response.StatusCode,
			"body":   string(raw),
		})
		return
	}

	var reasoning strings.Builder
	var text strings.Builder

	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer(
		make([]byte, 0, 64*1024),
		1024*1024,
	)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		log.Printf("[ask] %s", line)

		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimSpace(
			strings.TrimPrefix(line, "data:"),
		)

		if data == "[DONE]" {
			log.Println("[ask] stream closed by upstream")
			break
		}

		var chunk dto.ChatChunk

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Printf(
				"[ask] unparsable chunk skipped: %v",
				err,
			)
			continue
		}

		for _, choice := range chunk.Choices {
			if choice.Delta.ReasoningContent != "" {
				reasoning.WriteString(
					choice.Delta.ReasoningContent,
				)
			}

			if choice.Delta.Content != "" {
				text.WriteString(
					choice.Delta.Content,
				)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf(
			"[ask] stream read error: %v",
			err,
		)

		c.JSON(http.StatusBadGateway, gin.H{
			"error": "Stream interrupted",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ChatResponse{
		Model:     request.Model,
		Reasoning: reasoning.String(),
		Text:      text.String(),
	})
}