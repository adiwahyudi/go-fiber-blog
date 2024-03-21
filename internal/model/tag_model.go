package model

import "time"

type TagResponse struct {
	ID        uint           `json:"id,omitempty"`
	Name      string         `json:"name,omitempty"`
	Slug      string         `json:"slug,omitempty"`
	Posts     []PostResponse `json:"posts,omitempty"`
	CreatedAt *time.Time     `json:"created_at,omitempty"`
}

type CreateTagResponse struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name" validate:"required,max=20"`
}
