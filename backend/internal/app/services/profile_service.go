package services

import (
	"discord-backend/internal/app/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService struct {
	DB *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{DB: db}
}

func (p *ProfileService) GetProfileByID(profileID uuid.UUID) (*models.ProfileResponse, error) {
	var profile models.Profile
	if result := p.DB.Preload("Servers").First(&profile, profileID); result.Error != nil {
		return nil, result.Error
	}

	profileResponse := models.ProfileResponse{
		ID:        profile.ID,
		Name:      profile.Name,
		ImageURL:  profile.ImageURL,
		Email:     profile.Email,
		Servers:   profile.Servers,
		Members:   profile.Members,
		Channels:  profile.Channels,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	return &profileResponse, nil
}

func (p *ProfileService) UpdateProfile(userID uuid.UUID, updatedData models.Profile) (*models.Profile, error) {
	var profile models.Profile
	if result := p.DB.First(&profile, userID); result.Error != nil {
		return nil, result.Error
	}

	if result := p.DB.Model(&profile).Updates(updatedData); result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

func (p *ProfileService) CreateProfile(profile *models.Profile) (*models.ProfileResponse, error) {
	var existingProfile models.Profile
	if err := p.DB.Where("email = ?", profile.Email).Or("name = ?", profile.Name).First(&existingProfile).Error; err == nil {
		return nil, errors.New("email or username already taken")
	}

	hashedPassword, err := HashPassword(profile.Password)
	if err != nil {
		return nil, err
	}

	profile.Password = hashedPassword
	if err := p.DB.Create(profile).Error; err != nil {
		return nil, err
	}

	var profileResponse models.Profile
	if err := p.DB.Preload("Servers").Where("email = ?", profile.Email).First(&profileResponse).Error; err != nil {
		return nil, err
	}

	return &models.ProfileResponse{
		ID:        profileResponse.ID,
		Name:      profileResponse.Name,
		ImageURL:  profileResponse.ImageURL,
		Email:     profileResponse.Email,
		Servers:   profileResponse.Servers,
		Members:   profileResponse.Members,
		Channels:  profileResponse.Channels,
		CreatedAt: profileResponse.CreatedAt,
		UpdatedAt: profileResponse.UpdatedAt,
	}, nil
}

func (p *ProfileService) Authenticate(email, password string) (*models.ProfileResponse, error) {
	var profile models.Profile
	if err := p.DB.Preload("Servers").Where("email = ?", email).First(&profile).Error; err != nil {
		return nil, err
	}

	if !CheckPasswordHash(password, profile.Password) {
		return nil, errors.New("Invalid Credentials")
	}

	profileResponse := models.ProfileResponse{
		ID:        profile.ID,
		Name:      profile.Name,
		ImageURL:  profile.ImageURL,
		Email:     profile.Email,
		Servers:   profile.Servers,
		Members:   profile.Members,
		Channels:  profile.Channels,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}

	return &profileResponse, nil
}
