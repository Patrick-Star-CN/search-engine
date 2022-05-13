package searchService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetWebID(word string) (*models.DocID, error) {
	var docID *models.DocID

	result := database.DB.Where("word = ?", word).Find(docID)
	if result.Error != nil {
		return nil, result.Error
	}
	return docID, nil
}

func GetWebDoc(id int) (*models.DocRaw, error) {
	var docRaw *models.DocRaw

	result := database.DB.Where("id = ?", id).Find(docRaw)
	if result.Error != nil {
		return nil, result.Error
	}
	return docRaw, nil
}
