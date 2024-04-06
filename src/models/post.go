package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:char(36); primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255); not null" json:"title"`
	Content   string    `gorm:"type:text; not null" json:"content"`
	UserID    uuid.UUID `gorm:"type:char(36); not null" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()

	return nil
}
