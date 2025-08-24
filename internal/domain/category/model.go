package category

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name       string    `gorm:"uniqueIndex;size:100;not null"`
	CreatedAt  time.Time
	CreatedBy  uuid.UUID `gorm:"type:uuid"`
	ModifiedAt time.Time
	ModifiedBy uuid.UUID `gorm:"type:uuid"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
