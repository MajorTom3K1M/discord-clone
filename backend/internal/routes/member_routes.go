package routes

import (
	"discord-backend/internal/app/handlers"

	"github.com/gin-gonic/gin"
)

func MemberRoutes(protected *gin.RouterGroup, memberHandler *handlers.MemberHandler) {
	membersGroup := protected.Group("/members")
	{
		membersGroup.DELETE("/:memberId/servers/:serverId", memberHandler.KickMember)
		membersGroup.PATCH("/:memberId/servers/:serverId", memberHandler.UpdateMemberRole)
	}
}
