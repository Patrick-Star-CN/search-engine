package router

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/controllers/userController"
)

func userRouter(r *gin.RouterGroup) {
	r.POST("/submitCollection", userController.SubmitCollection)
	r.GET("/getCollection", userController.GetCollection)
	r.POST("/deleteCollection", userController.DelCollection)
}
