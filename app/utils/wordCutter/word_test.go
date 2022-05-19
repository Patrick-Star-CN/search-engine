package wordCutter

import (
	"fmt"
	"testing"
)

func TestWordCut(t *testing.T) {
	output := WordCut("小明硕士毕业于中国科学院计算所，后在日本京都大学深造")
	for _, word := range output {
		fmt.Printf("%s/", word)
	}
	fmt.Println()
}
