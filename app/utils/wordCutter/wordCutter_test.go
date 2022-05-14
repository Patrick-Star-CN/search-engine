package wordCutter

import (
	"fmt"
	"testing"
)

func TestWordCut(t *testing.T) {
	output := WordCut("今天的人民日报")
	for _, word := range output {
		fmt.Printf(" %s /", word)
	}
	fmt.Println()
}
