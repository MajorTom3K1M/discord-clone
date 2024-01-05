package services

import (
	"discord-backend/internal/app/models"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var jwtSecretKey = []byte(os.Getenv("SECRET_KEY"))

type TokenService struct {
	DB *gorm.DB
}

func NewTokenService(db *gorm.DB) *TokenService {
	return &TokenService{DB: db}
}

func GenerateTokens(profileID uuid.UUID, name string) (map[string]string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["profile_id"] = profileID
	claims["name"] = name
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	accessToken, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return nil, err
	}

	rtToken := jwt.New(jwt.SigningMethodHS256)

	rtClaims := rtToken.Claims.(jwt.MapClaims)

	rtClaims["profile_id"] = profileID
	rtClaims["name"] = name
	rtClaims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	refreshToken, err := rtToken.SignedString(jwtSecretKey)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid refresh token")
	}
}

func (t *TokenService) StoreRefreshToken(profileID uuid.UUID, refreshToken string, expiresAt time.Time) error {
	newRefreshToken := models.RefreshToken{
		ProfileID: profileID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}

	result := t.DB.Create(&newRefreshToken)
	return result.Error
}

func (t *TokenService) UpdateRefreshToken(profileID uuid.UUID, oldToken, newToken string) error {
	var refreshToken models.RefreshToken
	if err := t.DB.Where("profile_id = ? AND token = ?", profileID, oldToken).First(&refreshToken).Error; err != nil {
		return err
	}

	refreshToken.Token = newToken
	refreshToken.ExpiresAt = time.Now().Add(time.Hour * 72)

	return t.DB.Save(&refreshToken).Error
}

func (t *TokenService) FindRefreshToken(profileID uuid.UUID, refreshToken string) error {
	var refreshTokenModel models.RefreshToken
	if err := t.DB.Where("profile_id = ? AND token = ?", profileID, refreshToken).First(&refreshTokenModel).Error; err != nil {
		return err
	}
	return nil
}

func (t *TokenService) DeleteRefreshToken(profileID uuid.UUID, refreshToken string) error {
	return t.DB.Where("profile_id = ? AND token = ?", profileID, refreshToken).Delete(&models.RefreshToken{}).Error
}

func (t *TokenService) UpsertRefreshToken(profileID uuid.UUID, refreshToken string) error {
	newRefreshToken := models.RefreshToken{
		ProfileID: profileID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 72),
	}

	return t.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "profile_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"token", "expires_at"}),
	}).Create(&newRefreshToken).Error
}
