package router

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/controllers/searchController"
)

func searchRouterInit(r *gin.RouterGroup) {
	r.Any("/api/search", searchController.Search)
}
