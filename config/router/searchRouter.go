package router

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/controllers/cutDataController"
	"search-engine/app/controllers/searchController"
)

func searchRouterInit(r *gin.RouterGroup) {
	r.Any("/search", searchController.Search)
	r.GET("/cutData", cutDataController.CutData)
}
