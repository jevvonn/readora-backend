package entity

import (
	"time"
)

type Genre struct {
	Name  string `gorm:"primaryKey;type:varchar(255);not null;index" json:"name"`
	Books []Book `gorm:"many2many:book_genres;" json:"books"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
