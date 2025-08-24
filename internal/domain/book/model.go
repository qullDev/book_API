package book

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Book struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"size:200;not null"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null"`
	Description string    `gorm:"type:text"`
	ImageURL    string    `gorm:"type:text"`
	ReleaseYear int       `gorm:"not null"`
	Price       float64   `gorm:"not null"`
	TotalPage   int       `gorm:"not null"`
	Thickness   string    `gorm:"size:10;not null"`
	CreatedAt   time.Time
	CreatedBy   uuid.UUID `gorm:"type:uuid"`
	ModifiedAt  time.Time
	ModifiedBy  uuid.UUID `gorm:"type:uuid"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
