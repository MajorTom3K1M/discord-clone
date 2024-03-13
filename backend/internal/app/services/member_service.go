package services

import (
	"discord-backend/internal/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MemberService struct {
	DB *gorm.DB
}

func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{DB: db}
}

func (m *MemberService) UpdateMemberRole(serverID uuid.UUID, profileID uuid.UUID, memberID uuid.UUID, role models.MemberRole) (*models.Server, error) {
	var server models.Server

	err := m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Member{}).
			Where("id = ? AND profile_id <> ?", memberID, profileID).
			Update("role", role).Error; err != nil {
			return nil
		}

		if err := tx.Preload("Members", func(db *gorm.DB) *gorm.DB {
			return db.Order("members.role ASC").Preload("Profile")
		}).First(&server, "id = ? AND profile_id = ?", serverID, profileID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &server, nil
}
