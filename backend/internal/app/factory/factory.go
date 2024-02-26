package factory

import (
	"discord-backend/internal/app/handlers"
	"discord-backend/internal/app/services"

	"gorm.io/gorm"
)

type Factory struct {
	db *gorm.DB
}

func NewFactory(db *gorm.DB) *Factory {
	return &Factory{db: db}
}

func (f *Factory) NewProfileService() *services.ProfileService {
	return services.NewProfileService(f.db)
}

func (f *Factory) NewTokenService() *services.TokenService {
	return services.NewTokenService(f.db)
}

func (f *Factory) NewServerService() *services.ServerService {
	return services.NewServerService(f.db)
}

func (f *Factory) NewProfileHandler() *handlers.ProfileHandler {
	profileService := f.NewProfileService()
	return handlers.NewProfileHandler(profileService)
}

func (f *Factory) NewAuthHandler() *handlers.AuthHandler {
	profileService := f.NewProfileService()
	tokenService := f.NewTokenService()
	return handlers.NewAuthHandler(profileService, tokenService)
}

func (f *Factory) NewServerHandler() *handlers.ServerHandler {
	serverService := f.NewServerService()
	return handlers.NewServerHandler(serverService)
}
