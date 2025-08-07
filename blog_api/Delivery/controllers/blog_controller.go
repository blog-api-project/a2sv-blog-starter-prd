package controllers

import (
	"blog_api/Delivery/dtos"
	"blog_api/Delivery/utils"
	"blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type BlogController struct {
	blogUseCase usecases.IBlogUseCase
	imageUploader services.ImageUploader
}

func NewBlogController (blogUseCase usecases.IBlogUseCase,imageUploader services.ImageUploader) *BlogController{
	return &BlogController{
		blogUseCase : blogUseCase,
		imageUploader: imageUploader,
	}
}

func (bc *BlogController) CreateBlog(c *gin.Context) {
	var blogDTO dtos.BlogDto
	if err := c.ShouldBindWith(&blogDTO, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	var domainImages []models.UploadedImage
	for _, file := range form.File["images"] {
		f, err := file.Open()
		if err != nil {
			continue
		}
		defer f.Close()

		data, err := io.ReadAll(f)
		if err != nil {
			continue
		}

		image := models.UploadedImage{
			Filename: file.Filename,
			Size:     file.Size,
			Data:     data,
		}
		domainImages = append(domainImages, image)
	}

	imagePaths, err := bc.imageUploader.SaveImages(domainImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
		return
	}
	userID := c.GetString("userID")
	domainBlog := utils.ConvertToBlog(blogDTO, userID, imagePaths)


	if err := bc.blogUseCase.CreateBlog(domainBlog,userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog created successfully"})
}