package docRawService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetWebDoc(id int) (*models.DocRaw, error) {
	var docRaw *models.DocRaw

	result := database.DB.Where(
		&models.DocRaw{
			ID: id,
		}).Find(&docRaw)
	if result.Error != nil {
		return nil, result.Error
	}
	return docRaw, nil
}

func GetWebDocAll() ([]models.DocRaw, error) {
	var docRaws []models.DocRaw

	result := database.DB.Find(&docRaws)
	if result.Error != nil {
		return nil, result.Error
	}
	return docRaws, nil
}
