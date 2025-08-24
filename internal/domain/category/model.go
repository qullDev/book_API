package category

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name       string    `json:"name" gorm:"uniqueIndex;size:100;not null"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uuid.UUID `json:"created_by" gorm:"type:uuid"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy uuid.UUID `json:"modified_by" gorm:"type:uuid"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
