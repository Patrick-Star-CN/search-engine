package imgRawService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetImgRaw(id int) (*models.ImgRaw, error) {
	var imgRaw *models.ImgRaw

	result := database.DB.Where(
		&models.ImgRaw{
			ID: id,
		}).Find(&imgRaw)
	if result.Error != nil {
		return nil, result.Error
	}
	return imgRaw, nil
}

func GetImgRawAll() ([]models.ImgRaw, error) {
	var imgRaws []models.ImgRaw

	result := database.DB.Where("id > 185").Find(&imgRaws)
	if result.Error != nil {
		return nil, result.Error
	}
	return imgRaws, nil
}
