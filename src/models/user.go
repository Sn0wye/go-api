package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36); primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255); not null" json:"name"`
	Password  string    `gorm:"type:varchar(255); not null" json:"-"`
	Username  string    `gorm:"type:varchar(255); not null" json:"username"`
	Email     string    `gorm:"type:varchar(255); unique; not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()

	return nil
}
