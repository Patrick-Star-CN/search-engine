package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func JsonSuccessResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  "SUCCESS",
		"data": data,
	})
}
