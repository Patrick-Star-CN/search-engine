package models

type Collection struct {
	RawID int    `json:"id"`
	UID   int    `json:"uid"`
	URL   string `json:"url"`
	Title string `json:"title"`
}
