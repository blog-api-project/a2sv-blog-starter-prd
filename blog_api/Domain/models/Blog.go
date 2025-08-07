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