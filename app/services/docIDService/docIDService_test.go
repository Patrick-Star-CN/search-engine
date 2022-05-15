package docIDService

import "testing"

func TestGetWebID(t *testing.T) {
	output, err := GetWebID("a")
	if err != nil {
		println(err)
	}
	println((*output).ID)
}
