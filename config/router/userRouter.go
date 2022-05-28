package router

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/controllers/userController"
)

func userRouter(r *gin.RouterGroup) {
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)
	r.POST("/submitCollection", userController.SubmitCollection)
	r.GET("/getCollection", userController.GetCollection)
	r.DELETE("/deleteCollection", userController.DelCollection)
}
