package dto

type CreateMessageRequest struct {
	Message string `json:"message"`
}

type CreateMessageResponse struct {
	Message string `json:"message"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest mirrors the upstream payload in aichatendpoint.md.
// Stream is always forced to true on the wire.
type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

// ChatChunk is one SSE data: frame from the upstream.
// Delta carries reasoning_content before content appears.
type ChatChunk struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role             string `json:"role"`
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

type ChatResponse struct {
	Model     string `json:"model"`
	Reasoning string `json:"reasoning"`
	Text      string `json:"text"`
}