package services

import (
	"discord-backend/internal/app/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChannelService struct {
	DB *gorm.DB
}

func NewChannelService(db *gorm.DB) *ChannelService {
	return &ChannelService{DB: db}
}

func (c *ChannelService) CreateChannel(serverID, profileID uuid.UUID, name string, channelType models.ChannelType) (*models.Server, error) {
	var updatedServer models.Server

	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.Model(&models.Member{}).
			Where("server_id = ? AND profile_id = ? AND role IN ?",
				serverID, profileID, []string{string(models.Admin), string(models.Moderator)}).
			Count(&count).Error

		if err != nil {
			return err
		}

		if count == 0 {
			return errors.New("no matching members found")
		}

		channel := models.Channel{
			ProfileID: profileID,
			Name:      name,
			Type:      channelType,
			ServerID:  serverID,
		}

		if err := tx.Create(&channel).Error; err != nil {
			return err
		}

		if err := tx.Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("members.role ASC").Preload("Profile")
		}).Preload("Channels").
			First(&updatedServer, "id = ?", serverID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updatedServer, nil
}

func (c *ChannelService) DeleteChannel(serverID, profileID, channelID uuid.UUID) (*models.Server, error) {
	var updatedServer models.Server
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Member{}).
			Where("server_id = ? AND profile_id = ? AND role IN ?",
				serverID, profileID, []models.MemberRole{models.Admin, models.Moderator}).
			Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			return errors.New("no matching members found")
		}

		if err := tx.Where("id = ? AND server_id = ? AND name <> ?", channelID, serverID, "general").
			Delete(&models.Channel{}).Error; err != nil {
			return err
		}

		if err := tx.Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("members.role ASC").Preload("Profile")
		}).Preload("Channels").
			First(&updatedServer, "id = ?", serverID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updatedServer, nil
}

func (c *ChannelService) UpdateChannel(serverID, profileID, channelID uuid.UUID, updateData models.Channel) (*models.Server, error) {
	var updatedServer models.Server
	err := c.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Member{}).
			Where("server_id = ? AND profile_id = ? AND role IN ?",
				serverID, profileID, []models.MemberRole{models.Admin, models.Moderator}).
			Count(&count).Error; err != nil {
			return err
		}

		if count == 0 {
			return errors.New("no matching members found")
		}

		if err := tx.Where("id = ? AND server_id = ? AND name <> ?", channelID, serverID, "general").
			Updates(updateData).Error; err != nil {
			return err
		}

		if err := tx.Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("members.role ASC").Preload("Profile")
		}).Preload("Channels").
			First(&updatedServer, "id = ?", serverID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updatedServer, nil
}
