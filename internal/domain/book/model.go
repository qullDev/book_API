package book

import (
	"time"

	"github.com/google/uuid"
	"github.com/qullDev/book_API/internal/domain/category"
	"gorm.io/gorm"
)

type Book struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey"`
	Title       string            `gorm:"size:200;not null"`
	CategoryID  uuid.UUID         `gorm:"type:uuid;not null"`
	Category    category.Category `gorm:"foreignKey:CategoryID;references:ID"`
	Description string            `gorm:"type:text"`
	ImageURL    string            `gorm:"type:text"`
	ReleaseYear int               `gorm:"not null"`
	Price       float64           `gorm:"not null"`
	TotalPage   int               `gorm:"not null"`
	Thickness   string            `gorm:"size:10;not null"`
	CreatedAt   time.Time
	CreatedBy   uuid.UUID `gorm:"type:uuid"`
	ModifiedAt  time.Time
	ModifiedBy  uuid.UUID `gorm:"type:uuid"`
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	if b.TotalPage > 100 {
		b.Thickness = "tebal"
	} else {
		b.Thickness = "tipis"
	}
	b.ID = uuid.New()
	return
}
