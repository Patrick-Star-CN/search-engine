package searchService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetWebID(word string) (*models.DocID, error) {
	var docIDs *models.DocID

	result := database.DB.Where("word = ?", word).Find(docIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	return docIDs, nil
}
