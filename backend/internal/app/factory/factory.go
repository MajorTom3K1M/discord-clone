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

func (f *Factory) NewMemberService() *services.MemberService {
	return services.NewMemberService(f.db)
}

func (f *Factory) NewChannelService() *services.ChannelService {
	return services.NewChannelService(f.db)
}

func (f *Factory) NewConversationService() *services.ConversationService {
	return services.NewConversationService(f.db)
}

func (f *Factory) NewMessageService() *services.MessageService {
	return services.NewMessageService(f.db)
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

func (f *Factory) NewMemberHandler() *handlers.MemberHandler {
	memberService := f.NewMemberService()
	return handlers.NewMemberHandler(memberService)
}

func (f *Factory) NewChannelHandler() *handlers.ChannelHandler {
	channelService := f.NewChannelService()
	return handlers.NewChannelHandler(channelService)
}

func (f *Factory) NewConversationHandler() *handlers.ConversationHandler {
	conversationService := f.NewConversationService()
	return handlers.NewConversationHandler(conversationService)
}

func (f *Factory) NewWebsocketHandler() *handlers.WebsocketHandler {
	serverService := f.NewServerService()
	channelService := f.NewChannelService()
	messageService := f.NewMessageService()
	return handlers.NewWebsocketHandler(serverService, channelService, messageService)
}

func (f *Factory) NewMessageHandler() *handlers.MessageHandler {
	messageService := f.NewMessageService()
	return handlers.NewMessageHandler(messageService)
}
