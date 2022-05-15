package cutDataController

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/docIDService"
	"search-engine/app/services/docRawService"
	"search-engine/app/utils/wordCutter"
	"strconv"
)

func CutData(c *gin.Context) {
	docRaws, err := docRawService.GetWebDocAll()
	if err != nil {
		_ = c.AbortWithError(200, apiExpection.ServerError)
	}
	for _, v := range docRaws {
		words := wordCutter.WordCut(v.Title)
		for _, value := range words {
			docID, err := docIDService.GetWebID(value)
			if err != nil {
				_ = c.AbortWithError(200, apiExpection.ServerError)
			}
			if docID.Word != value {
				err := docIDService.CreateWebID(models.DocID{
					Word: value,
					ID:   strconv.Itoa(v.ID),
				})
				if err != nil {
					_ = c.AbortWithError(200, apiExpection.ServerError)
				}
			} else {
				err := docIDService.UpdateWebID(value, models.DocID{
					Word: value,
					ID:   docID.ID + ";" + strconv.Itoa(v.ID),
				})
				if err != nil {
					_ = c.AbortWithError(200, apiExpection.ServerError)
				}
			}
		}
	}
}
