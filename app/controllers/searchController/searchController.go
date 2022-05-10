package searchController

import (
	"github.com/gin-gonic/gin"
	"search-engine/app/apiExpection"
	"search-engine/app/services/searchService"
	"search-engine/app/utils/wordCutter"
)

func Search(c *gin.Context) {
	source := c.Query("word")
	wordsSlice := wordCutter.WordCut(source)

	for _, value := range wordsSlice {
		docs, err := searchService.GetWebID(value)
		if err != nil {
			_ = c.AbortWithError(200, apiExpection.ServerError)
		}

		for _, value := range docs.ID {

		}
	}

}
