package models

type DocRaw struct {
	id      int
	URL     string
	Title   string
	Context string
}

type DocRawScore struct {
	DocRaw
	Score int
}
