package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"` 
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string `json:"message"`
	AccessToken string `json:access_token`
}