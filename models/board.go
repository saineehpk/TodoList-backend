package models

import "gorm.io/gorm"

type Board struct {
	gorm.Model
	UserID uint   `gorm:"not null;index"`
	Name   string `gorm:"not null;size:100"`
	Cards  []Card `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
	Labels []Label `gorm:"foreignKey:BoardID;constraint:OnDelete:CASCADE"`
}
