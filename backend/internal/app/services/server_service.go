package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (s *ServerService) GetServerByProfileID(profileID uuid.UUID) (*models.Server, error) {
	var server models.Server
	err := s.DB.Joins("JOIN members ON members.server_id = servers.id").
		Where("members.profile_id = ?", profileID).
		First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) GetServer(profileID uuid.UUID, serverID uuid.UUID) (*models.Server, error) {
	var server models.Server

	err := s.DB.Preload("Members").Joins("JOIN members ON members.server_id = servers.id").
		Where("servers.id = ? AND members.profile_id = ?", serverID, profileID).
		First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) GetServerDetails(profileID uuid.UUID, serverID uuid.UUID) (*models.Server, error) {
	var server models.Server

	err := s.DB.Preload("Members", func(db *gorm.DB) *gorm.DB {
		return db.Order("members.role ASC").Preload("Profile")
	}).Preload("Channels", func(db *gorm.DB) *gorm.DB {
		return db.Order("channels.created_at ASC")
	}).Joins("JOIN members ON members.server_id = servers.id").
		Where("servers.id = ? AND members.profile_id = ?", serverID, profileID).First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) UpdateServerInviteCode(serverID uuid.UUID, profileID uuid.UUID) (*models.Server, error) {
	var server models.Server
	newInviteCode := uuid.New().String()

	result := s.DB.Model(&models.Server{}).Clauses(clause.Returning{}).
		Where("id = ? AND profile_id = ?", serverID, profileID).
		Update("invite_code", newInviteCode).Scan(&server)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &server, nil
}

func (s *ServerService) GetServerByInviteCode(inviteCode uuid.UUID, profileID uuid.UUID) (*models.Server, error) {
	var server models.Server

	err := s.DB.Joins("JOIN members ON members.server_id = servers.id").
		Where("servers.invite_code = ? AND members.profile_id = ?", inviteCode, profileID).
		First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) UpdateServerMember(inviteCode uuid.UUID, profileID uuid.UUID) (*models.Server, error) {
	var server models.Server

	if err := s.DB.Where("invite_code = ?", inviteCode).First(&server).Error; err != nil {
		return nil, err
	}

	member := models.Member{
		ProfileID: profileID,
		ServerID:  server.ID,
	}

	if err := s.DB.Create(&member).Error; err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) UpdateServer(profileID uuid.UUID, serverID uuid.UUID, name string, imageUrl string) (*models.Server, error) {
	var server models.Server

	updateData := models.Server{
		Name:     name,
		ImageURL: imageUrl,
	}

	result := s.DB.Model(&models.Server{}).Clauses(clause.Returning{}).
		Where("id = ? AND profile_id = ?", serverID, profileID).
		Updates(updateData).Scan(&server)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &server, nil
}

func (s *ServerService) LeaveServer(profileID, serverID uuid.UUID) (*models.Server, error) {

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var server models.Server
		// Step 1: Verify the server that have a this member and not admin
		if err := tx.Preload("Members", "profile_id = ?", profileID).
			Where("servers.id = ? AND servers.profile_id <> ?", serverID, profileID).
			First(&server).Error; err != nil {
			return err
		}

		// Step 2: Delete members with the specific profileId
		if err := tx.Where("server_id = ? AND profile_id = ?", serverID, profileID).
			Delete(&models.Member{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var updatedServer models.Server
	if err := s.DB.Preload("Members", func(db *gorm.DB) *gorm.DB {
		return db.Order("role ASC").Preload("Profile")
	}).First(&updatedServer, serverID).Error; err != nil {
		return nil, err
	}

	return &updatedServer, nil
}

func (s *ServerService) DeleteServer(profileID, serverID uuid.UUID) error {
	if err := s.DB.Where("id = ? AND profile_id = ?", serverID, profileID).
		Delete(&models.Server{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *ServerService) GetServerDefaultChannel(profileID, serverID uuid.UUID) (*models.Server, error) {
	var server models.Server

	err := s.DB.Preload("Channels", func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", "general").Order("created_at DESC")
	}).Joins("JOIN members ON members.server_id = servers.id").
		Where("servers.id = ? AND members.profile_id = ?", serverID, profileID).
		First(&server).Error

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerService) GetMember(serverID, profileID uuid.UUID) (*models.Member, error) {
	var member models.Member
	if err := s.DB.First(&member, "server_id = ? AND profile_id = ?", serverID, profileID).
		Error; err != nil {
		return nil, err
	}

	return &member, nil
}
