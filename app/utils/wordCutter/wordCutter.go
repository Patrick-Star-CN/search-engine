package wordCutter

import (
	"github.com/wangbin/jiebago"
)

var seg jiebago.Segmenter

func WordCut(source string) []string {
	Init()
	var wordsSlice []string
	wordMap := make(map[string]int)

	result := seg.CutForSearch(source, true)

	for {
		w, ok := <-result
		if !ok {
			break
		}
		_, found := wordMap[w]
		if !found {
			wordMap[w] = 1
		} else {
			wordMap[w]++
		}
	}

	for k, _ := range wordMap {
		wordsSlice = append(wordsSlice, k)
	}
	return wordsSlice
}

func Init() {
	seg.LoadDictionary("dict.txt")
}
