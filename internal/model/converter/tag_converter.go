package converter

import (
	"go-blog/internal/entity"
	"go-blog/internal/model"
)

func TagToResponse(tag *entity.Tag) *model.TagResponse {
	return &model.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
		//Posts:     tag.Posts,
		CreatedAt: tag.CreatedAt,
	}
}

func TagToCreateResponse(tag *entity.Tag) *model.TagResponse {
	return &model.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		Slug:      tag.Slug,
		CreatedAt: tag.CreatedAt,
	}
}

func TagToPostResponse(tag *entity.Tag) *model.TagResponse {
	return &model.TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
		Slug: tag.Slug,
	}
}
