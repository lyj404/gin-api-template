package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewOrgUnitRouter(orgHdlr *handler.OrgUnitHandler, group *gin.RouterGroup) {
	group.POST("/org-units", orgHdlr.CreateOrgUnit)
	group.PUT("/org-units/:id", orgHdlr.UpdateOrgUnit)
	group.DELETE("/org-units/:id", orgHdlr.DeleteOrgUnit)
	group.GET("/org-units/:id", orgHdlr.GetOrgUnit)
	group.GET("/org-units", orgHdlr.ListOrgUnits)
	group.GET("/org-units/tree", orgHdlr.GetOrgTree)
}
