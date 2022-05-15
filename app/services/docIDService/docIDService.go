package docIDService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetWebID(word string) (*models.DocID, error) {
	var docID *models.DocID

	result := database.DB.Where(
		&models.DocID{
			Word: word,
		}).Find(&docID)
	if result.Error != nil {
		return nil, result.Error
	}
	return docID, nil
}

func CreateWebID(docID models.DocID) error {
	result := database.DB.Create(&docID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateWebID(word string, docID models.DocID) error {
	result := database.DB.Model(models.DocID{}).Where(
		&models.DocID{
			Word: word,
		}).Updates(&docID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
