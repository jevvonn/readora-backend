package entity

import (
	"time"
)

type Genre struct {
	Name  string `gorm:"primaryKey;type:varchar(255);not null;index" json:"name,omitempty"`
	Books []Book `gorm:"many2many:book_genres;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"books,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
