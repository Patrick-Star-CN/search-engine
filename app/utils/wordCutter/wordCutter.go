package wordCutter

import (
	"github.com/yanyiwu/gojieba"
	"os"
	"path"
	"path/filepath"
)

func WordCut(source string) []string {
	dictDir := path.Join(filepath.Dir(os.Args[0]), "dict")
	jiebaPath := path.Join(dictDir, "jieba.dict.utf8")
	hmmPath := path.Join(dictDir, "hmm_model.utf8")
	userPath := path.Join(dictDir, "user.dict.utf8")
	idfPath := path.Join(dictDir, "idf.utf8")
	stopPath := path.Join(dictDir, "stop_words.utf8")
	x := gojieba.NewJieba(jiebaPath, hmmPath, userPath, idfPath, stopPath)

	defer x.Free()
	var wordsSlice []string
	wordMap := make(map[string]int)

	result := x.CutForSearch(source, false)

	for _, v := range result {
		_, found := wordMap[v]
		if !found {
			wordMap[v] = 1
		} else {
			wordMap[v]++
		}
	}

	for k, _ := range wordMap {
		wordsSlice = append(wordsSlice, k)
	}
	return wordsSlice
}
