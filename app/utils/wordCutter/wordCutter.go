package wordCutter

import "github.com/yanyiwu/gojieba"

func WordCut(source string) []string {
	var wordsSlice []string
	wordMap := make(map[string]int)
	x := gojieba.NewJieba()
	result := x.CutForSearch(source, true)

	for _, value := range result {
		_, found := wordMap[value]
		if !found {
			wordMap[value] = 1
		}
	}

	for k, _ := range wordMap {
		wordsSlice = append(wordsSlice, k)
	}
	return wordsSlice
}
