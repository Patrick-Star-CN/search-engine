package searchController

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
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

type SearchForm struct {
	Word     string `json:"word"`
	PaperNum int    `json:"paperNum"`
}

func Search(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var data Response
	var res SearchForm

	errBind := c.ShouldBindJSON(&res)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	word, wordErr := url.QueryUnescape(res.Word)
	var wordsSlice []string
	if wordErr != nil {
		wordsSlice = wordCutter.WordCut(res.Word)
	} else {
		wordsSlice = wordCutter.WordCut(word)
	}
	docs := make(map[int]int)

	for _, value := range wordsSlice {
		docID, err := docIDService.GetWebID(value)
		if err != nil {
			log.Println("table web_id error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
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

	for i := 10 * (res.PaperNum - 1); i < res.PaperNum*10 && i < len(docs_); i++ {
		doc, err := docRawService.GetWebDoc(docs_[i].id)
		if err != nil {
			log.Println("table web_doc error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
		data.Data = append(data.Data, *doc)
	}
	data.Length = len(docs_)
	utils.JsonSuccessResponse(c, data)
}
