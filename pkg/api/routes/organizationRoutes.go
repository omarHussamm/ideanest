package routes

import (
	"demo/pkg/api/handlers"

	"github.com/gin-gonic/gin"
)

func OrganizationRoute(router *gin.RouterGroup) {
	router.POST("/", handlers.CreateOrganization)
	router.GET("/", handlers.GetAllOrganizations)
	router.GET("/:organizationId", handlers.GetOrganization)
	router.PUT("/:organizationId", handlers.UpdateOrganization)
	router.DELETE("/:organizationId", handlers.DeleteOrganization)
	router.POST("/:organizationId/invite", handlers.InviteUsertoOrganization)
}
