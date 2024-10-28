package repo

import (
	"errors"
	"fmt"

	"github.com/codedbyshoe/goxclone/internal/models"
	"gorm.io/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *PostRepo) GetUsersPosts(userId uint) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("User").
		Preload("Comments").
		Preload("Likes").
		Preload("Reposts").
		Where("user_id = ?", userId).
		Find(&posts)
	return posts, result.Error
}

func (r *PostRepo) GetPost(postId uint) (models.Post, error) {
	var post models.Post
	result := r.db.Preload("User").Where("id = ?", postId).First(&post)

	return post, result.Error
}

func (r *PostRepo) UpdatePost(postId uint, user_id uint, content string) (*models.Post, error) {
	post := &models.Post{}
	err := r.db.Where("id = ? AND user_id = ?", postId, user_id).First(post).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Post not found or unauthorized")
		}
		return nil, err
	}

	updates := map[string]interface{}{
		"content": content,
	}

	err = r.db.Model(post).Updates(updates).Save(post).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Preload("User").
		Preload("Comments").
		Preload("Likes").
		Preload("Reposts").
		First(post, postId).Error

	return post, err
}

func (r *PostRepo) DeletePost(postId uint, userId uint) error {
	result := r.db.Where("id = ? AND user_id = ?", postId, userId).Delete(&models.Post{})

	if result.RowsAffected == 0 {
		return fmt.Errorf("Post not found or you are not authorized")
	}

	return result.Error
}
