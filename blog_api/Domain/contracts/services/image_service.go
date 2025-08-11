package services

import (
	"blog_api/Domain/models"
)

type ImageUploader interface {
	SaveImages(files []models.UploadedImage) ([]string, error)
}