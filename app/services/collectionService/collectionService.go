package collectionService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetCollection(uid int, url string) (*models.Collection, error) {
	var collection *models.Collection

	result := database.DB.Where(
		&models.Collection{
			UID: uid,
			URL: url,
		}).Find(&collection)
	if result.Error != nil {
		return nil, result.Error
	}
	return collection, nil
}

func GetCollectionAll(uid int) ([]models.Collection, error) {
	var collections []models.Collection

	result := database.DB.Where(
		&models.Collection{
			UID: uid,
		}).Find(&collections)
	if result.Error != nil {
		return nil, result.Error
	}
	return collections, nil
}

func CreateCollection(data models.Collection) error {
	result := database.DB.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
