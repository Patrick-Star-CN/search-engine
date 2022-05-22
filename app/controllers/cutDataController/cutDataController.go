package cutDataController

import (
	"github.com/gin-gonic/gin"
	"log"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/docIDService"
	"search-engine/app/services/docRawService"
	"search-engine/app/utils/wordCutter"
	"strconv"
)

func CutData(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	docRaws, err := docRawService.GetWebDocAll()
	if err != nil {
		log.Println("table web_doc error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	for _, v := range docRaws {
		words := wordCutter.WordCut(v.Title)
		for _, value := range words {
			docID, err := docIDService.GetWebID(value)
			if err != nil {
				log.Println("table web_id error")
				_ = c.AbortWithError(200, apiExpection.ServerError)
				return
			}
			if docID.Word != value {
				err := docIDService.CreateWebID(models.DocID{
					Word: value,
					ID:   strconv.Itoa(v.ID),
				})
				if err != nil {
					log.Println("table web_id error")
					_ = c.AbortWithError(200, apiExpection.ServerError)
					return
				}
			} else {
				err := docIDService.UpdateWebID(value, models.DocID{
					Word: value,
					ID:   docID.ID + ";" + strconv.Itoa(v.ID),
				})
				if err != nil {
					log.Println("table web_id error")
					_ = c.AbortWithError(200, apiExpection.ServerError)
					return
				}
			}
		}
	}
}
