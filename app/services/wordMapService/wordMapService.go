package wordMapService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func CreatMap(data models.WordMap) error {
	result := database.DB.Create(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetMap(preWord string) (*models.WordMap, error) {
	var wordMap *models.WordMap

	result := database.DB.Where(
		&models.WordMap{
			PreWord: preWord,
		}).Find(&wordMap)
	if result.Error != nil {
		return nil, result.Error
	}
	return wordMap, nil
}

func UpdateMap(data models.WordMap) error {
	result := database.DB.Model(models.WordMap{}).Where(
		&models.WordMap{
			PreWord: data.PreWord,
		}).Updates(&data)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
