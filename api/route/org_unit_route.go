package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
	"github.com/lyj404/gin-api-template/api/middleware"
)

func NewOrgUnitRouter(orgHdlr *handler.OrgUnitHandler, rbac *middleware.RBACMiddleware, group *gin.RouterGroup) {
	group.POST("/org-units", orgHdlr.CreateOrgUnit)
	group.PUT("/org-units/:id", orgHdlr.UpdateOrgUnit)
	group.DELETE("/org-units/:id", orgHdlr.DeleteOrgUnit)
	group.GET("/org-units/:id", rbac.CheckPermission("org:manage"), orgHdlr.GetOrgUnit)
	group.GET("/org-units", rbac.CheckPermission("org:manage"), orgHdlr.ListOrgUnits)
	group.GET("/org-units/tree", rbac.CheckPermission("org:manage"), orgHdlr.GetOrgTree)
}
