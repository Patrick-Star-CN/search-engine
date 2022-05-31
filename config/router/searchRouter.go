package router

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/controllers/cutDataController"
	"search-engine/app/controllers/searchController"
)

func searchRouterInit(r *gin.RouterGroup) {
	r.POST("/search", searchController.Search)
	r.POST("/searchImg", searchController.SearchImg)
	r.GET("/cutData", cutDataController.CutData)
	r.GET("/cutImgData", cutDataController.CutImgData)
	r.POST("/submitHistory", searchController.SubmitHistory)
	r.GET("/getRelated", searchController.GetHistory)
}
