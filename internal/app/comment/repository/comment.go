package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type CommentPostgreSQLItf interface {
	CreateComment(req entity.Comment) error
	GetComments(filter dto.GetCommentsQuery) ([]entity.Comment, error)
	GetSpecificComment(req entity.Comment) (entity.Comment, error)
	DeleteComment(req entity.Comment) error
}

type CommentPostgreSQL struct {
	db     *gorm.DB
	logger logger.LoggerItf
}

func NewCommentPostgreSQL(db *gorm.DB, logger logger.LoggerItf) CommentPostgreSQLItf {
	return &CommentPostgreSQL{db, logger}
}

func (r *CommentPostgreSQL) CreateComment(req entity.Comment) error {
	log := "[CommentPostgreSQL][CreateComment]"

	err := r.db.Create(&req).Error
	if err != nil {
		r.logger.Error(log, err)
		return err
	}

	return nil
}

func (r *CommentPostgreSQL) GetComments(filter dto.GetCommentsQuery) ([]entity.Comment, error) {
	log := "[CommentPostgreSQL][GetComments]"
	var comments []entity.Comment
	query := r.db.Preload("User").Model(&entity.Comment{})

	if filter.BookId != "" {
		query = query.Where("book_id = ?", filter.BookId)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Page > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if filter.TopCommentId != "" {
		query = query.Order("CASE WHEN id = '" + filter.TopCommentId + "' THEN 0 ELSE 1 END")
	}

	if filter.SortBy != "" {
		query = query.Order(filter.SortBy + " " + filter.SortOrder)
	}

	err := query.Find(&comments).Error
	if err != nil {
		r.logger.Error(log, err)
		return []entity.Comment{}, err
	}

	return comments, nil
}

func (r *CommentPostgreSQL) GetSpecificComment(req entity.Comment) (entity.Comment, error) {
	log := "[CommentPostgreSQL][GetSpecificComment]"
	var comment entity.Comment
	err := r.db.Where(&req).First(&comment).Error

	if err != nil {
		r.logger.Error(log, err)
		return entity.Comment{}, err
	}

	return comment, nil
}

func (r *CommentPostgreSQL) DeleteComment(req entity.Comment) error {
	log := "[CommentPostgreSQL][DeleteComment]"
	err := r.db.Delete(&req).Error

	if err != nil {
		r.logger.Error(log, err)
		return err
	}

	return nil
}
