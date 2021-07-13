package repositories

import (
	"devbook-api/src/models"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Posts struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *Posts {
	var post Posts
	db.Model(&post).Association("Users")

	return &Posts{db}
}

func (postRepository Posts) Create(post models.Post) (models.Post, error) {
	postId := uuid.NewV4()
	newPost := models.Post{
		ID:       postId,
		Title:    post.Title,
		Content:  post.Content,
		AuthorID: post.AuthorID,
	}
	result := postRepository.db.Omit("likes").Create(&newPost)
	if result.Error != nil {
		return models.Post{}, result.Error
	}

	createdPost := postRepository.db.Preload("User").Find(&newPost, "id = ?", postId)
	if createdPost.Error != nil {
		return models.Post{}, result.Error
	}
	newPost.Sanitize()
	return newPost, nil
}

func (postRepository Posts) FindPostById(postId string) (models.Post, error) {
	var post models.Post
	result := postRepository.db.Preload("User").Find(&post, "id = ?", postId)

	if result.Error != nil {
		return models.Post{}, result.Error
	}
	post.Sanitize()
	return post, nil
}

func (postRepository Posts) GetPosts(userId string) ([]models.Post, error) {
	var posts []models.Post

	result := postRepository.db.Joins("User").Find(&posts).Joins("join followers on followers.user_id = posts.author_id").Where("followers.follower_id = ? ", userId).Scan(&posts)
	if result.Error != nil {
		return []models.Post{}, result.Error
	}
	for _, post := range posts {
		post.Sanitize()
	}
	return posts, nil
}