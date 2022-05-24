package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/collectionService"
	"search-engine/app/utils"
)

func SubmitCollection(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req models.Collection

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	collection, err := collectionService.GetCollection(req.UID, req.URL)
	if err != nil {
		log.Println("table collection error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if collection.URL == req.URL {
		utils.JsonSuccessResponse(c, "REPEATED", nil)
		return
	}

	e := collectionService.CreateCollection(req)
	if e != nil {
		log.Println("table collection error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
