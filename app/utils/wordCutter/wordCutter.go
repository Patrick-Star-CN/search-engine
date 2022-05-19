package wordCutter

import "github.com/yanyiwu/gojieba"

func WordCut(source string) []string {
	x := gojieba.NewJieba()
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
