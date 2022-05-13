package searchController

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/searchService"
	"search-engine/app/utils"
	"search-engine/app/utils/wordCutter"
	"sort"
	"strconv"
	"strings"
)

type Response struct {
	data   []models.DocRaw
	length int
}

func Search(c *gin.Context) {
	var data Response
	source := c.Query("word")
	wordsSlice := wordCutter.WordCut(source)
	docs := make(map[int]models.DocRawScore)
	var docs_ []models.DocRawScore

	for _, value := range wordsSlice {
		docID, err := searchService.GetWebID(value)
		if err != nil {
			_ = c.AbortWithError(200, apiExpection.ServerError)
		}
		IDs := strings.Split(docID.ID, ";")
		for _, value := range IDs {
			id, _ := strconv.Atoi(value)
			v, found := docs[id]
			if !found {
				doc, err := searchService.GetWebDoc(id)
				if err != nil {
					_ = c.AbortWithError(200, apiExpection.ServerError)
				}
				docScore := models.DocRawScore{DocRaw: *doc, Score: 1}
				docs[id] = docScore
			} else {
				v.Score++
			}
		}
	}

	for _, value := range docs {
		docs_ = append(docs_, value)
	}
	sort.SliceStable(docs_, func(i, j int) bool {
		return docs_[i].Score > docs_[j].Score
	})
	for _, value := range docs_ {
		data.data = append(data.data, value.DocRaw)
	}
	data.length = len(docs_)
	utils.JsonSuccessResponse(c, data)
}
