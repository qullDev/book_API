package book

import (
	"time"

	"github.com/google/uuid"
	"github.com/qullDev/book_API/internal/domain/category"
	"gorm.io/gorm"
)

type Book struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primaryKey"`
	Title       string            `json:"title" gorm:"size:200;not null"`
	CategoryID  uuid.UUID         `json:"category_id" gorm:"type:uuid;not null"`
	Category    category.Category `json:"category" gorm:"foreignKey:CategoryID;references:ID"`
	Description string            `json:"description" gorm:"type:text"`
	ImageURL    string            `json:"image_url" gorm:"type:text"`
	ReleaseYear int               `json:"release_year" gorm:"not null"`
	Price       float64           `json:"price" gorm:"not null"`
	TotalPage   int               `json:"total_page" gorm:"not null"`
	Thickness   string            `json:"thickness" gorm:"size:10;not null"`
	CreatedAt   time.Time         `json:"created_at"`
	CreatedBy   uuid.UUID         `json:"created_by" gorm:"type:uuid"`
	ModifiedAt  time.Time         `json:"modified_at"`
	ModifiedBy  uuid.UUID         `json:"modified_by" gorm:"type:uuid"`
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
