package controller

import (
	"fmt"
	"log"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"time"
	"context"

	"auth/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"crypto/sha256"
	"encoding/hex"
)

type AuthController struct {
	redis *redis.Client
}

func NewAuthController(redis *redis.Client) *AuthController {
	return &AuthController{
		redis: redis,
	}
}

func (ac *AuthController) Hello (c *gin.Context) {	
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Hello Too",
	})
}

func (ac *AuthController) Login (c *gin.Context) {	
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

	//in case have cookie but still hitting login
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		err = ac.RevokeRefreshToken(refreshToken)
		if err != nil {
			log.Printf("Failed to revoke refresh token: %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to revoke refresh token",
			})
			return
		}
	}

	accessToken, err := GenerateAccessToken(request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	refreshToken, err = GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	err = ac.StoreRefreshToken(refreshToken, request.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store refresh token",
		})
		return
	}	

	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*7, // 7 hari
		"/",
		"",
		false,  // Secure
		true,  // HttpOnly
	)	

	response := dto.LoginResponse{
		Message: "Valid Credentials",
		AccessToken: accessToken,
	}

	c.JSON(http.StatusAccepted, response)
}

func (ac *AuthController) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		log.Printf("Failed to get refresh_token cookie: %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token cookie not found",
		})
		return
	}

	userId, err := ac.GetUserIdFromRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Failed to get refresh_token cookie: %v", err)

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	err = ac.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Failed to revoke refresh token: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to revoke refresh token",
		})
		return
	}

	accessToken, err := GenerateAccessToken(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	refreshToken, err = GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}
	
	err = ac.StoreRefreshToken(refreshToken, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to store refresh token",
		})
		return
	}	

	c.SetCookie(
		"refresh_token",
		refreshToken,
		60*60*24*7, // 7 hari
		"/",
		"",
		false,  // Secure
		true,  // HttpOnly
	)	

	response := dto.LoginResponse{
		Message: "Refresh token rotated",
		AccessToken: accessToken,
	}

	c.JSON(http.StatusAccepted, response)	
}

func (ac *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Refresh token cookie not found",
		})
		return
	}

	err = ac.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Failed to revoke refresh token: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to logout",
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
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
		"exp": time.Now().Add(5 * time.Minute).Unix(),
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

func HashRefreshToken(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken))

	return hex.EncodeToString(hash[:])
}

func (ac *AuthController) StoreRefreshToken(refreshToken string, userId string) error {
	tokenHash := HashRefreshToken(refreshToken)

	return ac.redis.Set(
		context.Background(),
		"refresh_token:"+tokenHash,
		userId,
		7*24*time.Hour,
	).Err()
}

func (ac *AuthController) GetUserIdFromRefreshToken(refreshToken string) (string, error) {
	tokenHash := HashRefreshToken(refreshToken)

	userId, err := ac.redis.Get(
		context.Background(),
		"refresh_token:"+tokenHash,
	).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("refresh token not found")
	}

	if err != nil {
		return "", err
	}

	return userId, nil
}

func (ac *AuthController) RevokeRefreshToken(refreshToken string) error {
	tokenHash := HashRefreshToken(refreshToken)

	return ac.redis.Del(
		context.Background(),
		"refresh_token:"+tokenHash,
	).Err()
}