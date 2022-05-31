package imgIDService

import (
	"search-engine/app/models"
	"search-engine/config/database"
)

func GetImgID(word string) (*models.ImgID, error) {
	var imgID *models.ImgID

	result := database.DB.Where(
		&models.ImgID{
			Word: word,
		}).Find(&imgID)
	if result.Error != nil {
		return nil, result.Error
	}
	return imgID, nil
}

func CreateImgID(imgID models.ImgID) error {
	result := database.DB.Create(&imgID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func UpdateImgID(word string, imgID models.ImgID) error {
	result := database.DB.Model(models.ImgID{}).Where(
		&models.ImgID{
			Word: word,
		}).Updates(&imgID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
