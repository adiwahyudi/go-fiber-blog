package converter

import (
	"go-blog/internal/entity"
	"go-blog/internal/model"
)

func PostToResponse(post *entity.Post) *model.PostResponse {
	var tags []*model.TagResponse
	for _, tag := range post.Tags {
		tags = append(tags, TagToPostResponse(tag))
	}
	return &model.PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Slug:    post.Slug,
		Content: post.Content,
		Tags:    tags,
		User: model.UserOnPost{
			ID:       post.User.ID,
			Name:     post.User.Name,
			Username: post.User.Username,
		},
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}
}
