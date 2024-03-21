package model

import "time"

type PostResponse struct {
	ID          uint           `json:"id,omitempty"`
	Title       string         `json:"name,omitempty"`
	Slug        string         `json:"slug,omitempty"`
	Content     string         `json:"content,omitempty"`
	Tags        []*TagResponse `json:"tags,omitempty"`
	User        UserOnPost     `json:"user,omitempty"`
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
}

type CreatePostRequest struct {
	Title   string              `json:"title" validate:"required,max=100"`
	Content string              `json:"content" validate:"required"`
	Tags    []CreateTagResponse `json:"tags"`
}

type UserOnPost struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type SearchPostRequest struct {
	Username string     `json:"username" form:"username" validate:"min=1,max=30"`
	Sort     string     `json:"sort" form:"sort" validate:"min=1"`
	Title    string     `json:"title" form:"title" validate:"max=100"`
	Tags     []string   `json:"tags" form:"tags"`
	UserId   string     `json:"-"`
	Paginate Pagination `json:"paginate"`
}
