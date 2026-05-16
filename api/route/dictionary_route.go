package route

import (
	"github.com/gin-gonic/gin"
	"github.com/lyj404/gin-api-template/api/handler"
)

func NewDictionaryRouter(h *handler.DictionaryHandler, group *gin.RouterGroup) {
	dict := group.Group("/dict")
	{
		dict.GET("", h.ListDict)
		dict.GET("/:id", h.GetDict)
		dict.POST("", h.CreateDict)
		dict.PUT("/:id", h.UpdateDict)
		dict.DELETE("/:id", h.DeleteDict)

		dict.GET("/:id/details", h.ListDictDetails)
		dict.POST("/:id/details", h.CreateDictDetail)
		dict.PUT("/:id/details/:detailId", h.UpdateDictDetail)
		dict.DELETE("/:id/details/:detailId", h.DeleteDictDetail)
	}
}

func NewPublicDictionaryRouter(h *handler.DictionaryHandler, group *gin.RouterGroup) {
	group.GET("/dict-info/:type", h.GetDictInfoByType)
}
