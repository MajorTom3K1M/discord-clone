package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServerService struct {
	DB *gorm.DB
}

func NewServerService(db *gorm.DB) *ServerService {
	return &ServerService{DB: db}
}

func (s *ServerService) CreateServer(profileID uuid.UUID, name string, imageUrl string) (*models.Server, error) {
	inviteCode := uuid.New().String()

	server := models.Server{
		ProfileID:  profileID,
		Name:       name,
		ImageURL:   imageUrl,
		InviteCode: inviteCode,
		Channels: []models.Channel{
			{Name: "general", ProfileID: profileID},
		},
		Members: []models.Member{
			{ProfileID: profileID, Role: models.Admin},
		},
	}

	tx := s.DB.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&server).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) GetServers(profileID uuid.UUID) ([]models.Server, error) {
	var servers []models.Server
	err := s.DB.Joins("JOIN members ON members.server_id = servers.id").
		Where("members.profile_id = ?", profileID).
		Distinct("servers.*").
		Find(&servers).Error

	if err != nil {
		return nil, err
	}

	return servers, nil
}

func (s *ServerService) GetServer(profileID uuid.UUID, serverID uuid.UUID) (*models.Server, error) {
	var server models.Server

	err := s.DB.Joins("JOIN members ON members.server_id = servers.id").
		Where("servers.id = ? AND members.profile_id = ?", serverID, profileID).
		First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}
