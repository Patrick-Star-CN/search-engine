package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/collectionService"
	"search-engine/app/services/sessionService"
	"search-engine/app/utils"
)

type CollectionForm struct {
	ID  []int `json:"id"`
	UID int   `json:"uid"`
}

func SubmitCollection(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req models.Collection

	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		log.Println("session error:" + errSession.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	collection, err := collectionService.GetCollection(user.ID, req.URL)
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

func GetCollection(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		log.Println("session error:" + errSession.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	collections, err := collectionService.GetCollectionAll(user.ID)
	if err != nil {
		log.Println("table collection error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", gin.H{
		"length": len(collections),
		"data":   collections,
	})
}

func DelCollection(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req CollectionForm

	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		log.Println("session error:" + errSession.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	for _, i := range req.ID {
		err := collectionService.DelCollection(i, user.ID)
		if err != nil {
			log.Println("table collection error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
