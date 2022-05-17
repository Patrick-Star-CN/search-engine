package searchController

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/docIDService"
	"search-engine/app/services/docRawService"
	"search-engine/app/utils"
	"search-engine/app/utils/wordCutter"
	"sort"
	"strconv"
	"strings"
)

type Response struct {
	Data   []models.DocRaw
	Length int
}

type DocIDs struct {
	id    int
	score int
}

func Search(c *gin.Context) {
	var data Response
	source := c.Query("word")
	wordsSlice := wordCutter.WordCut(source)
	docs := make(map[int]int)

	for _, value := range wordsSlice {
		docID, err := docIDService.GetWebID(value)
		if err != nil {
			_ = c.AbortWithError(200, apiExpection.ServerError)
		}
		if docID.Word != value {
			continue
		}

		IDs := strings.Split(docID.ID, ";")
		for _, value := range IDs {
			id, _ := strconv.Atoi(value)
			_, found := docs[id]
			if !found {
				docs[id] = 1
			} else {
				docs[id]++
			}
		}
	}

	docs_ := make([]DocIDs, len(docs))
	index := 0
	for key, value := range docs {
		docs_[index].id = key
		docs_[index].score = value
		index++
	}
	sort.SliceStable(docs_, func(i, j int) bool {
		return docs_[i].score > docs_[j].score
	})
	for i := 0; i < 50 && i < len(docs_); i++ {
		doc, err := docRawService.GetWebDoc(docs_[i].id)
		if err != nil {
			_ = c.AbortWithError(200, apiExpection.ServerError)
		}
		data.Data = append(data.Data, *doc)
	}
	data.Length = len(docs_)
	utils.JsonSuccessResponse(c, data)
}
