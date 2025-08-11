// infrastructure/image_service.go
package infrastructure

import (
	"blog_api/Domain/models"
	"fmt"
	"os"
	"path/filepath"
)

type ImageService struct {
	UploadDir string
}

func NewImageService(uploadDir string) *ImageService {
	return &ImageService{UploadDir: uploadDir}
}

func (s *ImageService) SaveImages(images []models.UploadedImage) ([]string, error) {
	var savedPaths []string

	for _, img := range images {
		fileName := generateFileName(img.Filename)
		fullPath := filepath.Join(s.UploadDir, fileName)

		if err := os.WriteFile(fullPath, img.Data, 0644); err != nil {
			return nil, fmt.Errorf("failed to save image %s: %w", img.Filename, err)
		}
		savedPaths = append(savedPaths, fullPath)
	}

	return savedPaths, nil
}

func generateFileName(original string) string {
	ext := filepath.Ext(original)
	name := filepath.Base(original)
	return name  + ext
}
