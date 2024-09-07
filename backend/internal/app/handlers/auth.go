package handlers

import (
	"discord-backend/internal/app/models"
	"discord-backend/internal/app/services"
	"discord-backend/internal/app/utils"
	customErrors "discord-backend/internal/app/utils"
	"errors"
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
	// var profile models.Profile
	var profile struct {
		Name     string `json:"name"`
		ImageURL string `json:"imageUrl"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newProfile := models.Profile{
		Name:     profile.Name,
		ImageURL: profile.ImageURL,
		Email:    profile.Email,
		Password: profile.Password,
	}

	currentProfile, err := h.ProfileService.CreateProfile(&newProfile)
	if err != nil {
		if errors.Is(err, customErrors.ErrEmailOrUsernameTaken) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already taken"})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokens, err := services.GenerateTokens(currentProfile.ID, currentProfile.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT token"})
		return
	}

	if err := h.TokenService.UpsertRefreshToken(currentProfile.ID, tokens["refreshToken"]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upsert JWT token to database"})
		return
	}

	host := c.Request.Host
	domain := utils.ExtractBaseDomain(host)
	maxAge := int(time.Hour * 72 / time.Second)
	c.SetCookie("refresh_token", tokens["refreshToken"], maxAge, "/", domain, true, true)
	c.SetCookie("access_token", tokens["accessToken"], maxAge, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful", "profile": currentProfile})
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

	// c.Header("Authorization", "Bearer "+tokens["accessToken"])

	host := c.Request.Host
	domain := utils.ExtractBaseDomain(host)
	maxAge := int(time.Hour * 72 / time.Second)
	c.SetCookie("refresh_token", tokens["refreshToken"], maxAge, "/", domain, true, true)
	c.SetCookie("access_token", tokens["accessToken"], maxAge, "/", domain, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "profile": profile})
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

	host := c.Request.Host
	domain := utils.ExtractBaseDomain(host)
	c.SetCookie("access_token", "", -1, "/", domain, false, true)
	c.SetCookie("refresh_token", "", -1, "/", domain, true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	tx := h.TokenService.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction start error"})
		return
	}

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

	if err := h.TokenService.UpsertRefreshToken(profileID, tokens["refreshToken"]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error upsert JWT token to database"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit error"})
		return
	}

	maxAge := int(time.Hour * 72 / time.Second)
	host := c.Request.Host
	domain := utils.ExtractBaseDomain(host)
	c.SetCookie("refresh_token", tokens["refreshToken"], maxAge, "/", domain, true, true)
	c.SetCookie("access_token", tokens["accessToken"], maxAge, "/", domain, false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed"})
}

func (h *AuthHandler) ServerRefresh(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Tokens refreshed", "access_token": tokens})
}
