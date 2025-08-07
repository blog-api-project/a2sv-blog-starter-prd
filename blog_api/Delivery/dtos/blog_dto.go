package dtos


type BlogDto struct {
  Title   string   `form:"title"`
  Content string   `form:"content"`
  Tags    []string `form:"tags"`
}

type BlogQueryDto struct {
    Page       int      `form:"page"`
    PageSize    int      `form:"page_size"`
    SortBy     string   `form:"sort_by"`
  Title      string  `form:"title"`
  Author     string  `form:"author"`
  Tags       string  `form:"tags"`
}

type PaginationMetadataDTO struct{
	TotalPages    int  `json:"total_pages"`
	CurrentPage   int  `json:"current_page"`
	TotalPosts    int   `json:"total_posts"`
	PageSize      int   `json:"total_posts"`
}
