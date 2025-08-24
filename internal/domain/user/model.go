package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Username   string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password   string    `json:"password" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	CreatedBy  uuid.UUID `json:"created_by" gorm:"type:uuid"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy uuid.UUID `json:"modified_by" gorm:"type:uuid"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
