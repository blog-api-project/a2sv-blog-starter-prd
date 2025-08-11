package models

import "time"

type Blog struct {
	ID            string
	AuthorID      string
	Title         string
	Content       string
	ImageURL      []string
	Tags          []string
	PostedAt      time.Time
	LikeCount     int
	DislikeCount  int
	CommentCount  int
	ShareCount    int
	AISuggestion  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UploadedImage struct {
	Filename string
	Size     int64
	Data     []byte
}

type BlogQuery struct {
    Page     int
    PageSize int
    SortBy   string
	Title    string
	Author   string
	Tags     string


}
type PaginationMeta struct {
	TotalPages   int
	CurrentPage  int
	TotalPosts   int
	PostsPerPage int
}
type UserBlogInteraction struct{
	ID          string
	UserID      string
	BlogID      string
	Action      string
	CreatedAt time.Time

}


