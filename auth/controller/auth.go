package controller

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"
	"auth/dto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Hello (c *gin.Context) {	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Hello Too",
	})
}

func Login (c *gin.Context) {	
	var request dto.LoginRequest
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if request.Username != "admin" || request.Password != "password"{
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid Credentials",
		})
		return
	}

	accessToken, err := GenerateAccessToken(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*7, // 7 hari
		"/",
		"",
		true,  // Secure
		true,  // HttpOnly
	)	

	response := dto.LoginResponse{
		Message: "Valid Credentials",
		AccessToken: accessToken,
	}

	c.JSON(http.StatusAccepted, response)
}

func GenerateAccessToken(username string) (string, error) {
	privateKeyData, err := os.ReadFile("key/private.pem")
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(privateKey)
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}