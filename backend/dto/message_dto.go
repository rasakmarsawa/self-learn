package dto

type CreateMessageRequest struct {
	Message string `json:"message"`
}

type CreateMessageResponse struct {
	Message string `json:"message"`
}