package searchController

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/url"
	"search-engine/app/apiExpection"
	"search-engine/app/models"
	"search-engine/app/services/docIDService"
	"search-engine/app/services/docRawService"
	"search-engine/app/services/imgIDService"
	"search-engine/app/services/imgRawService"
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

type ResponseImg struct {
	Data   []models.ImgRaw
	Length int
}

type Relation struct {
	word string
	num  int
}

type IDs struct {
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

	// 绑定请求到结构体
	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	// 过滤词在“-”后面，切词处理
	words := strings.Split(req.Word, "-")

	// 判断是否需要解码，需要就解，并且调用jieba进行切词
	word, wordErr := url.QueryUnescape(words[0])
	var wordsSlice []string
	if wordErr != nil {
		wordsSlice = wordCutter.WordCut(req.Word)
	} else {
		wordsSlice = wordCutter.WordCut(word)
	}
	docs := make(map[int]int)

	// 切出来的每一个词都调用映射表获取对应的id，且统计各个id出现的次数以此判断关联性
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

	// 将map转换为切片数组，然后按出现次数进行排序
	docs_ := make([]IDs, len(docs))
	index := 0
	for key, value := range docs {
		docs_[index].id = key
		docs_[index].score = value
		index++
	}
	sort.SliceStable(docs_, func(i, j int) bool {
		return docs_[i].score > docs_[j].score
	})

	// 获取页码对应的十条数据
	i := 10 * (req.PaperNum - 1)
	for j := 1; j < 10 && i < len(docs_); j++ {
		doc, err := docRawService.GetWebDoc(docs_[i].id)
		if err != nil {
			log.Println("table web_doc error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
		i++
		if (len(words) > 1 && !strings.ContainsAny(doc.Title, words[1])) || len(words) == 1 {
			// 如果过滤词且这条搜索结果中有过滤词就跳过
			data.Data = append(data.Data, *doc)
		}
	}
	data.Length = len(docs_)
	utils.JsonSuccessResponse(c, "SUCCESS", data)
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

	// 去除过滤词，查询映射表中是否已有对应数据
	req.PreWord = strings.Split(req.PreWord, "-")[0]
	req.Word = strings.Split(req.Word, "-")[0]
	wordMap, e := wordMapService.GetMap(req.PreWord)
	if e != nil {
		log.Println("table word_map error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if wordMap.PreWord == req.PreWord {
		// 映射表中已存在前一次搜索词的映射
		var num string
		flag := true
		words := strings.Split(wordMap.Word, ";")
		for i, word := range words {
			// 判断前搜索词与现搜索词的映射是否已存在，是则给该映射权重+1
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
			// 映射不存在则新加一条映射
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
		// 前搜索词不存在映射，新建映射
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

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func GetHistory(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	// 获取参数并转码
	word_ := c.Query("word")
	word, e := url.QueryUnescape(word_)
	if e != nil {
		word = word_
	}

	// 查询是否存在本次搜索词对应的映射，不存在则返回空
	wordMap, err := wordMapService.GetMap(word)
	if err != nil {
		log.Println("table word_map error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if wordMap.PreWord != word {
		utils.JsonSuccessResponse(c, "SUCCESS", nil)
		return
	}

	// 将映射切片并排序
	words := strings.Split(wordMap.Word, ";")
	nums := strings.Split(wordMap.Num, ";")
	res := make([]Relation, len(words))
	for i, s := range words {
		res[i].num, _ = strconv.Atoi(nums[i])
		res[i].word = s
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].num > res[j].num
	})
	for i, re := range res {
		words[i] = re.word
	}
	utils.JsonSuccessResponse(c, "SUCCESS", words)
}

func SearchImg(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var data ResponseImg
	var req SearchForm

	// 绑定请求到结构体
	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	// 过滤词在“-”后面，切词处理
	words := strings.Split(req.Word, "-")

	// 判断是否需要解码，需要就解，并且调用jieba进行切词
	word, wordErr := url.QueryUnescape(words[0])
	var wordsSlice []string
	if wordErr != nil {
		wordsSlice = wordCutter.WordCut(req.Word)
	} else {
		wordsSlice = wordCutter.WordCut(word)
	}
	imgs := make(map[int]int)

	// 切出来的每一个词都调用映射表获取对应的id，且统计各个id出现的次数以此判断关联性
	for _, value := range wordsSlice {
		imgID, err := imgIDService.GetImgID(value)
		if err != nil {
			log.Println("table img_id error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
		if imgID.Word != value {
			continue
		}

		IDs := strings.Split(imgID.ID, ";")
		for _, value := range IDs {
			id, _ := strconv.Atoi(value)
			_, found := imgs[id]
			if !found {
				imgs[id] = 1
			} else {
				imgs[id]++
			}
		}
	}

	// 将map转换为切片数组，然后按出现次数进行排序
	imgs_ := make([]IDs, len(imgs))
	index := 0
	for key, value := range imgs {
		imgs_[index].id = key
		imgs_[index].score = value
		index++
	}
	sort.SliceStable(imgs_, func(i, j int) bool {
		return imgs_[i].score > imgs_[j].score
	})

	// 获取页码对应的十条数据
	i := 10 * (req.PaperNum - 1)
	for j := 1; j < 10 && i < len(imgs_); j++ {
		img, err := imgRawService.GetImgRaw(imgs_[i].id)
		if err != nil {
			log.Println("table img_doc error")
			_ = c.AbortWithError(200, apiExpection.ServerError)
			return
		}
		i++
		if (len(words) > 1 && !strings.ContainsAny(img.Title, words[1])) || len(words) == 1 {
			// 如果过滤词且这条搜索结果中有过滤词就跳过
			data.Data = append(data.Data, *img)
		}
	}
	data.Length = len(imgs_)
	utils.JsonSuccessResponse(c, "SUCCESS", data)
}
