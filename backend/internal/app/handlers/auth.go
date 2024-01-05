package handlers

import (
	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	ProfileService *services.ProfileService
	TokenService   *services.TokenService
}

func NewAuthHandler(profileService *services.ProfileService, tokenService *services.TokenService) *AuthHandler {
	return &AuthHandler{
		ProfileService: profileService,
		TokenService:   tokenService,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var profile models.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.ProfileService.CreateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.ProfileService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokens, err := services.GenerateTokens(profile.ID, profile.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT token"})
		return
	}

	if err := h.TokenService.UpsertRefreshToken(profile.ID, tokens["refreshToken"]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upsert JWT token to database"})
		return
	}

	maxAge := int(time.Hour * 72 / time.Second)
	c.SetCookie("refresh_token", tokens["refreshToken"], maxAge, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "profile": profile, "accessToken": tokens["accessToken"]})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	var signOutRequest struct {
		ProfileID string `json:"profileId"`
	}

	if err := c.ShouldBindJSON(&signOutRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profileID, err := uuid.Parse(signOutRequest.ProfileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UUID from claims"})
		return
	}

	refreshTokenCookie, err := c.Request.Cookie("refresh_token")
	if err == nil {
		h.TokenService.DeleteRefreshToken(profileID, refreshTokenCookie.Value)
	}

	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenCookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing refresh token"})
		return
	}

	claims, err := services.VerifyToken(refreshTokenCookie.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	profileIDString := claims["profile_id"].(string)
	name := claims["name"].(string)
	profileID, err := uuid.Parse(profileIDString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse UUID from claims"})
		return
	}

	if err := h.TokenService.FindRefreshToken(profileID, refreshTokenCookie.Value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token"})
		return
	}

	tokens, err := services.GenerateTokens(profileID, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.Header("Authorization", "Bearer "+tokens["accessToken"])

	maxAge := int(time.Hour * 72 / time.Second)
	c.SetCookie("refresh_token", tokens["refreshToken"], maxAge, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}
