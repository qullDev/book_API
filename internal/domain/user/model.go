package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	Username   string    `gorm:"uniqueIndex;size:50;not null"`
	Password   string    `gorm:"not null"`
	CreatedAt  time.Time
	CreatedBy  uuid.UUID `gorm:"type:uuid"`
	ModifiedAt time.Time
	ModifiedBy uuid.UUID `gorm:"type:uuid"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
