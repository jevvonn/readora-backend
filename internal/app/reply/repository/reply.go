package repository

import (
	"github.com/jevvonn/readora-backend/internal/domain/dto"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type ReplyPostgreSQLItf interface {
	Create(req entity.Reply) error
	GetSpecificReply(req entity.Reply) (entity.Reply, error)
	GetRepliesByCommentId(commentId string, filter dto.GetRepliesQuery) ([]entity.Reply, error)
}

type ReplyPostgreSQL struct {
	db  *gorm.DB
	log logger.LoggerItf
}

func NewReplyPostgreSQL(db *gorm.DB, log logger.LoggerItf) ReplyPostgreSQLItf {
	return &ReplyPostgreSQL{db, log}
}

func (r *ReplyPostgreSQL) Create(req entity.Reply) error {
	log := "[ReplyPostgreSQL][Create]"

	err := r.db.Create(&req).Error
	if err != nil {
		r.log.Error(log, err)
	}

	return nil
}

func (r *ReplyPostgreSQL) GetSpecificReply(req entity.Reply) (entity.Reply, error) {
	log := "[ReplyPostgreSQL][GetSpecificReply]"

	var reply entity.Reply
	err := r.db.Where(&req).First(&reply).Error
	if err != nil {
		r.log.Error(log, err)
		return reply, err
	}

	return reply, nil
}

func (r *ReplyPostgreSQL) GetRepliesByCommentId(commentId string, filter dto.GetRepliesQuery) ([]entity.Reply, error) {
	log := "[ReplyPostgreSQL][GetRepliesByCommentId]"

	var replies []entity.Reply
	query := r.db.Debug().Preload("User").Model(&entity.Reply{}).Where("comment_id = ?", commentId)

	if filter.Limit >= 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Page >= 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if filter.ParentId != "" {
		query = query.Where("parent_id = ?", filter.ParentId)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	query = query.Order("created_at " + filter.SortOrder)

	err := query.Find(&replies).Error
	if err != nil {
		r.log.Error(log, err)
		return nil, err
	}

	return replies, nil
}
