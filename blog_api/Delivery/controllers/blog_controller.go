package controllers

import (
	"blog_api/Delivery/dtos"
	"blog_api/Delivery/utils"
	"blog_api/Domain/contracts/services"
	usecases "blog_api/Domain/contracts/usecases"
	"blog_api/Domain/models"
	"io"

	"math"

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

func (bc *BlogController) GetBlogs (c *gin.Context){

	var BlogQueryParams dtos.BlogQueryDto
	if err := c.ShouldBindQuery(&BlogQueryParams); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid Request"})
		return
	}

	domainQuery := utils.ConvertToBlogQuery(BlogQueryParams)

	blogs,total,err := bc.blogUseCase.GetBlogs(domainQuery)
	if err != nil{
		c.JSON(http.StatusBadGateway,gin.H{"error":err.Error()})
		return 
	}
	paginationMeta := dtos.PaginationMetadataDTO{
		TotalPages:  int(math.Ceil(float64(total) / float64(domainQuery.PageSize))),
		CurrentPage: domainQuery.Page,
		TotalPosts:  total,
		PageSize:    domainQuery.PageSize,
}

	c.JSON(http.StatusOK,gin.H{
		"blog":blogs,
		"pagination":paginationMeta,

	})


}

func (ct *BlogController) UpdateBlogHandler(c *gin.Context) {
	blogID := c.Param("id")
	userID := c.GetString("userID")

	var updatedBlogDTO dtos.BlogDto
	if err := c.ShouldBindWith(&updatedBlogDTO, binding.FormMultipart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	blogModel := utils.ConvertToBlog(updatedBlogDTO,userID,nil)

	updatedBlog, err := ct.blogUseCase.UpdateBlog(blogModel, blogID, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Blog updated successfully",
		"updated_blog": updatedBlog,
	})
}

func (bc *BlogController) DeleteBlogHandler(c *gin.Context){
	blogID := c.Param("id")
	userID := c.GetString("userID")

	err := bc.blogUseCase.DeleteBlog(blogID,userID)
	if err != nil{
		c.JSON(http.StatusBadGateway,gin.H{
			"error":err.Error()})
			return
	}
	c.JSON(http.StatusOK,gin.H{"Message":"Blog Deleted Successfully"})
	

}