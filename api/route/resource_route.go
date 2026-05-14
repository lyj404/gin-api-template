package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewResourceRouter(resourceHdlr *handler.ResourceHandler, group *gin.RouterGroup) {
	group.POST("/resources", resourceHdlr.CreateResource)
	group.PUT("/resources/:id", resourceHdlr.UpdateResource)
	group.DELETE("/resources/:id", resourceHdlr.DeleteResource)
	group.GET("/resources/:id", resourceHdlr.GetResource)
	group.GET("/resources", resourceHdlr.ListResources)
}
