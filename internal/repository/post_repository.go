package repository

import (
	"github.com/sirupsen/logrus"
	"go-blog/internal/entity"
	"go-blog/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepository struct {
	Repository[entity.Post]
	Log *logrus.Logger
}

func NewPostRepository(log *logrus.Logger) *PostRepository {
	return &PostRepository{
		Log: log,
	}
}

func (r *Repository[T]) Find(db *gorm.DB, request *model.SearchPostRequest) ([]entity.Post, int64, error) {
	var posts []entity.Post

	err := db.
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Slug")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Username")
		}).
		Scopes(r.filterPostScopes(request)).
		Offset((request.Paginate.Page - 1) * request.Paginate.Size).
		Limit(request.Paginate.Size).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	err = db.
		Model(&entity.Post{}).
		Scopes(r.filterPostScopes(request)).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (r *Repository[T]) FindBySlug(db *gorm.DB, entity *T, slug string) error {
	return db.
		Where("slug = ?", slug).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Slug")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("ID", "Name", "Username")
		}).Take(entity).Error
}
func (r *Repository[T]) filterPostScopes(request *model.SearchPostRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if username := request.Username; username != "" {
			tx = tx.
				Joins("inner join users u on u.id = posts.user_id").
				Where("u.username = ?", username)
		}
		if len(request.Tags) > 0 && request.Tags[0] != "" {
			tx = tx.
				Joins("inner join post_tags pt on pt.post_id = posts.id").
				Joins("inner join tags t on pt.tag_id = t.id").
				Where("t.slug IN ?", request.Tags)
		}

		if title := request.Title; title != "" {
			title = "%" + title + "%"
			tx = tx.Where("title LIKE ?", title)
		}

		if sort := request.Sort; sort != "" {
			switch sort {
			case "latest":
				tx = tx.Order("created_at desc")
			case "oldest":
				tx = tx.Order("created_at asc")
			default:
				tx = tx.Order("created_at desc")
			}
		}

		return tx

	}
}
func (r *Repository[T]) filterPostClauses(request *model.SearchPostRequest) []clause.Expression {
	clauses := make([]clause.Expression, 0)
	if title := request.Title; title != "" {
		title = "%" + title + "%"
		clauses = append(clauses, clause.Like{
			Column: "title",
			Value:  title,
		})
	}

	if sort := request.Sort; sort != "" {
		var sortFilter clause.OrderBy
		switch sort {
		case "latest":
			sortFilter = clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column: clause.Column{Name: "created_at"},
						Desc:   true,
					},
				},
			}
		case "oldest":
			sortFilter = clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column: clause.Column{Name: "created_at"},
						Desc:   false,
					},
				},
			}
		default:
			sortFilter = clause.OrderBy{
				Columns: []clause.OrderByColumn{
					{
						Column: clause.Column{Name: "published_at"},
						Desc:   true,
					},
				},
			}
		}
		clauses = append(clauses, sortFilter)
	}
	return clauses

}
