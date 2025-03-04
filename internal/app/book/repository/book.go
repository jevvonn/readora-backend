package repository

import (
	"github.com/google/uuid"
	"github.com/jevvonn/readora-backend/internal/constant"
	"github.com/jevvonn/readora-backend/internal/domain/entity"
	"github.com/jevvonn/readora-backend/internal/infra/logger"
	"gorm.io/gorm"
)

type BookPostgreSQLItf interface {
	Create(req entity.Book) error
	GetBooks(filter GetBooksFilter) ([]entity.Book, error)
}

type BookPostgreSQL struct {
	db  *gorm.DB
	log logger.LoggerItf
}

type GetBooksFilter struct {
	Search    string
	Genre     string
	Limit     int
	Page      int
	SortBy    string
	SortOrder string
	OwnerID   uuid.UUID
	Role      string
}

func NewBookPostgreSQL(db *gorm.DB, log logger.LoggerItf) BookPostgreSQLItf {
	return &BookPostgreSQL{db, log}
}

func (r *BookPostgreSQL) Create(req entity.Book) error {
	err := r.db.Create(&req).Error
	if err != nil {
		r.log.Error("[BookPostgreSQL][Create]", err)
	}

	return err
}

func (r *BookPostgreSQL) GetBooks(filter GetBooksFilter) ([]entity.Book, error) {
	var books []entity.Book
	query := r.db.Model(&entity.Book{}).Preload("Owner").Preload("Genres")

	if filter.Role == constant.RoleAdmin {
		query = query.Joins("JOIN users ON users.id = books.owner_id")
		query = query.Where("users.role = ?", constant.RoleAdmin)
	}

	if filter.OwnerID != uuid.Nil {
		query = query.Where("owner_id = ?", filter.OwnerID.String())
	}

	if filter.Search != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Search+"%").Or("author ILIKE ?", "%"+filter.Search+"%").Or("description ILIKE ?", "%"+filter.Search+"%")
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Page > 0 {
		query = query.Offset((filter.Page - 1) * filter.Limit)
	}

	if filter.SortBy != "" {
		query = query.Order(filter.SortBy + " " + filter.SortOrder)
	}

	if filter.Genre != "" {
		query = query.Joins("JOIN book_genres ON book_genres.book_id = books.id").Joins("JOIN genres ON genres.name = book_genres.genre_name")
		query = query.Where("genres.name ILIKE ?", filter.Genre)
	}

	err := query.Find(&books).Error
	if err != nil {
		r.log.Error("[BookPostgreSQL][GetBooks]", err)
		return nil, err
	}

	return books, nil
}
