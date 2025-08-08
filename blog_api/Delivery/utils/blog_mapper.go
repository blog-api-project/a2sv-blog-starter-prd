package utils

import (
  "blog_api/Delivery/dtos"
  "blog_api/Domain/models"
  "time"
)

func ConvertToBlog(dto dtos.BlogDto, authorID string, imageURLs []string) *models.Blog {
 return &models.Blog{
   Title:     dto.Title,
   Content:   dto.Content,
   Tags:      dto.Tags,
   AuthorID:  authorID,
   ImageURL: imageURLs,
   PostedAt:  time.Now(),
   CreatedAt: time.Now(),
   UpdatedAt: time.Now(),
 }
}


func ConvertToBlogQuery (dto dtos.BlogQueryDto) *models.BlogQuery{
	return &models.BlogQuery{
		Page: dto.Page,
		PageSize: dto.PageSize,
		SortBy: dto.SortBy,
		Title : dto.Title,
		Author : dto.Author,
		Tags : dto.Tags,

	}
}

