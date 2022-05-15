package models

type DocRaw struct {
	ID      int
	URL     string
	Title   string
	Content string
}

type DocRawScore struct {
	DocRaw
	Score int
}
