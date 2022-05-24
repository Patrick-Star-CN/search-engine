package searchController

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/docIDService"
	"search-engine/app/services/docRawService"
	"search-engine/app/services/wordMapService"
	"search-engine/app/utils"
	"search-engine/app/utils/wordCutter"
	"sort"
	"strconv"
	"strings"
)

type ResponseDoc struct {
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

type HistoryForm struct {
	PreWord string `json:"preWord"`
	Word    string `json:"word"`
}

func Search(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var data ResponseDoc
	var req SearchForm

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	words := strings.Split(req.Word, "-")

	word, wordErr := url.QueryUnescape(words[0])
	var wordsSlice []string
	if wordErr != nil {
		wordsSlice = wordCutter.WordCut(req.Word)
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

	if len(words) > 1 {
		docID, err := docIDService.GetWebID(words[1])
		if err != nil {
			log.Println("table web_id error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
		if docID.Word == words[1] {
			IDs := strings.Split(docID.ID, ";")
			for _, value := range IDs {
				id, _ := strconv.Atoi(value)
				_, found := docs[id]
				if found {
					delete(docs, id)
				}
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

	for i := 10 * (req.PaperNum - 1); i < req.PaperNum*10 && i < len(docs_); i++ {
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

func SubmitHistory(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req HistoryForm
	var data models.WordMap

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	req.PreWord = strings.Split(req.PreWord, "-")[0]
	req.Word = strings.Split(req.Word, "-")[0]
	wordMap, e := wordMapService.GetMap(req.PreWord)
	if e != nil {
		log.Println("table word_map error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if wordMap.PreWord == req.PreWord {
		var num string
		flag := true
		words := strings.Split(wordMap.Word, ";")
		for i, word := range words {
			if word == req.Word {
				nums := strings.Split(wordMap.Num, ";")
				num_, _ := strconv.Atoi(nums[i])
				nums[i] = strconv.Itoa(num_ + 1)
				num = strings.Join(nums, ";")
				flag = false
				break
			}
		}
		if !flag {
			data.Num = num
			data.PreWord = wordMap.PreWord
			data.Word = wordMap.Word
		} else {
			data.Num = wordMap.Num + ";1"
			data.PreWord = wordMap.PreWord
			data.Word = wordMap.Word + ";" + req.Word
		}
		err := wordMapService.UpdateMap(data)
		if err != nil {
			log.Println("table word_map error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
	} else {
		data.Num = "1"
		data.PreWord = req.PreWord
		data.Word = req.Word
		err := wordMapService.CreatMap(data)
		if err != nil {
			log.Println("table word_map error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
	}

	utils.JsonSuccessResponse(c, "")
}
